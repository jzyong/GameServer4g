package util

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync/atomic"
)

//@
type Uint32IdAllocator struct {
	id uint32
}

//@
func NewUint32IdAllocator() *Uint32IdAllocator {
	return &Uint32IdAllocator{}
}

//@
func (this *Uint32IdAllocator) Get() uint32 {
	id := atomic.AddUint32(&this.id, 1)
	if id == 0 {
		id = atomic.AddUint32(&this.id, 1)
	}
	return id
}

//@
func GenerateSessionId() (string, error) {
	k := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return "", nil
	}
	return hex.EncodeToString(k), nil
}
