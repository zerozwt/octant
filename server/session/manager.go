package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"
)

type Manager interface {
	GenerateSessionKey() string
	Set(sessKey string, data any, ttl int64)
	Get(sessKey string) (any, bool)
	Del(sessKey string)
}

func GetManager() Manager { return mgr }

type manager struct {
	sync.Mutex

	enc  cipher.Block
	incr atomic.Int64
	data map[string]*unit
}

var mgr *manager

func init() {
	mgr = &manager{
		data: map[string]*unit{},
	}

	var raw [16]byte
	rand.Read(raw[:])
	mgr.enc, _ = aes.NewCipher(raw[:])

	go mgr.checkExpiration()
}

type unit struct {
	data any

	ttl      int64
	expireAt int64
}

func (u *unit) expired(ts time.Time) bool {
	return u.expireAt <= ts.Unix()
}

func (u *unit) refresh(ts time.Time) {
	u.expireAt = ts.Unix() + u.ttl
}

func (m *manager) GenerateSessionKey() string {
	var buf [24]byte

	binary.LittleEndian.PutUint64(buf[:], uint64(time.Now().UnixMicro()))
	binary.BigEndian.PutUint64(buf[8:], uint64(time.Now().Unix()))
	binary.BigEndian.PutUint64(buf[16:], uint64(m.incr.Add(1)))

	m.enc.Encrypt(buf[8:], buf[8:])

	return base64.RawURLEncoding.EncodeToString(buf[:])
}

func (m *manager) Set(sessKey string, data any, ttl int64) {
	if ttl <= 0 {
		return
	}

	m.Lock()
	defer m.Unlock()

	m.data[sessKey] = &unit{
		data:     data,
		ttl:      ttl,
		expireAt: time.Now().Unix() + ttl,
	}
}

func (m *manager) Get(sessKey string) (any, bool) {
	m.Lock()
	defer m.Unlock()

	if unit, ok := m.data[sessKey]; ok {
		now := time.Now()
		if unit.expired(now) {
			delete(m.data, sessKey)
			return nil, false
		}
		unit.refresh(now)
		return unit.data, true
	}

	return nil, false
}

func (m *manager) Del(sessKey string) {
	m.Lock()
	defer m.Unlock()

	delete(m.data, sessKey)
}

func (m *manager) checkExpiration() {
	time.Sleep(time.Second)

	n := 0
	max := 512

	m.Lock()
	defer m.Unlock()

	key2del := make([]string, 0, max)
	now := time.Now()

	for key, value := range m.data {
		if value.expired(now) {
			key2del = append(key2del, key)
		}
		n += 1
		if n >= max {
			break
		}
	}

	for _, key := range key2del {
		delete(m.data, key)
	}

	go m.checkExpiration()
}

func getSessionData[T any](key string) (ret T, ok bool) {
	var raw any
	raw, ok = GetManager().Get(key)

	if !ok {
		return
	}

	ret, ok = raw.(T)
	return
}
