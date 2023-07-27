package bridge

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/zerozwt/etcdutil"
	"github.com/zerozwt/swe"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func CreateEtcdBridge(client *clientv3.Client) Bridge {
	ret := &etcdBridge{
		client: client,
		rooms:  map[int64]bool{},
		shadow: map[int64]bool{},
	}
	ret.watcher = etcdutil.NewWatcher(client, ret, etcdRoomPrefix, clientv3.WithPrefix())
	return ret
}

var etcdRoomPrefix string = "room_"

type etcdBridge struct {
	client  *clientv3.Client
	watcher *etcdutil.Watcher

	lock   sync.Mutex
	rooms  map[int64]bool
	shadow map[int64]bool
	recv   Receiver

	inReset     atomic.Bool
	recvChanged atomic.Bool

	stopped int32
}

func (b *etcdBridge) roomKey(roomID int64) string {
	return etcdRoomPrefix + fmt.Sprint(roomID)
}

func (b *etcdBridge) AddRoom(roomID int64) error {
	_, err := b.client.KV.Put(context.Background(), b.roomKey(roomID), "-")
	return err
}

func (b *etcdBridge) DelRoom(roomID int64) error {
	_, err := b.client.KV.Delete(context.Background(), b.roomKey(roomID))
	return err
}

func (b *etcdBridge) SetReceiver(recv Receiver) {
	if recv == nil {
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.recv = recv
	if !b.inReset.Load() {
		for roomID := range b.rooms {
			b.recv.OnAddRoom(roomID)
		}
	} else {
		b.recvChanged.Store(true)
	}
}

func (b *etcdBridge) Start() error {
	b.watcher.Start()
	return nil
}

func (b *etcdBridge) Stop() error {
	atomic.AddInt32(&b.stopped, 1)
	return b.watcher.Stop()
}

func (b *etcdBridge) key2room(key string) (int64, error) {
	if !strings.HasPrefix(key, etcdRoomPrefix) {
		return 0, fmt.Errorf("invalid room key %s", key)
	}

	ret, err := strconv.ParseInt(key[len(etcdRoomPrefix):], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse key %s failed: %v", key, err)
	}

	return ret, nil
}

func (b *etcdBridge) OnPut(key, value []byte) {
	roomID, err := b.key2room(string(key))
	if err != nil {
		return
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.rooms[roomID]; !ok {
		b.rooms[roomID] = true
		if b.recv != nil && !b.inReset.Load() {
			b.recv.OnAddRoom(roomID)
		}
	}
}

func (b *etcdBridge) OnDelete(key []byte) {
	roomID, err := b.key2room(string(key))
	if err != nil {
		return
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.rooms[roomID]; ok {
		delete(b.rooms, roomID)
		if b.recv != nil && !b.inReset.Load() {
			b.recv.OnDelRoom(roomID)
		}
	}
}

func (b *etcdBridge) OnError(err error) {
	if atomic.LoadInt32(&b.stopped) > 0 {
		return
	}

	swe.CtxLogger(nil).Error("etcd watcher error %v, try reconnect ...", err)

	b.watcher = etcdutil.NewWatcher(b.client, b, etcdRoomPrefix, clientv3.WithPrefix())
	b.watcher.Start()
}

func (b *etcdBridge) OnResetBegin() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.inReset.Store(true)

	b.shadow = b.rooms
	b.rooms = map[int64]bool{}
}

func (b *etcdBridge) OnResetEnd() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.inReset.Store(false)

	if b.recvChanged.Load() {
		for roomID := range b.rooms {
			b.recv.OnAddRoom(roomID)
		}
		b.recvChanged.Store(false)
		b.shadow = map[int64]bool{}
		return
	}

	toDel := keySubstract(b.shadow, b.rooms)
	toAdd := keySubstract(b.rooms, b.shadow)

	if b.recv != nil {
		for _, roomID := range toDel {
			b.recv.OnDelRoom(roomID)
		}

		for _, roomID := range toAdd {
			b.recv.OnAddRoom(roomID)
		}
	}

	b.shadow = map[int64]bool{}
}

func keySubstract[K comparable, V any](a, b map[K]V) []K {
	ret := []K{}
	for k := range a {
		if _, ok := b[k]; !ok {
			ret = append(ret, k)
		}
	}
	return ret
}
