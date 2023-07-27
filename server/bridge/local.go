package bridge

import "sync"

func CreateLocalBridge() Bridge {
	return &localBridge{
		rooms: map[int64]bool{},
	}
}

type localBridge struct {
	lock  sync.Mutex
	rooms map[int64]bool

	recv Receiver
}

func (b *localBridge) AddRoom(roomID int64) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.rooms[roomID]; !ok {
		b.rooms[roomID] = true
		if b.recv != nil {
			b.recv.OnAddRoom(roomID)
		}
	}

	return nil
}

func (b *localBridge) DelRoom(roomID int64) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.rooms[roomID]; ok {
		delete(b.rooms, roomID)
		if b.recv != nil {
			b.recv.OnDelRoom(roomID)
		}
	}

	return nil
}

func (b *localBridge) SetReceiver(recv Receiver) {
	if recv == nil {
		return
	}
	b.lock.Lock()
	defer b.lock.Unlock()

	b.recv = recv
	for roomID := range b.rooms {
		b.recv.OnAddRoom(roomID)
	}
}

func (b *localBridge) Start() error { return nil }
func (b *localBridge) Stop() error  { return nil }
