package utils

import (
	"crypto/rand"
	"encoding/binary"
	"sync/atomic"
	"time"
)

const (
	highMask = (1 << 43) - 1
	lowMask  = (1 << 20) - 1
	intMask  = (1 << 31) - 1
)

var idLow atomic.Int32

func init() {
	var buf [4]byte
	rand.Read(buf[:])
	n := binary.BigEndian.Uint32(buf[:])
	idLow.Store(int32(n & intMask))
}

func GenerateID() int64 {
	return ((time.Now().Unix() & highMask) << 20) | int64(idLow.Add(1)&lowMask)
}
