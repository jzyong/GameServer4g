package util

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync/atomic"
	"time"
)

//
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

/*
组成：0(1 bit) | timestamp in milli second (41 bit) | machine id (10 bit) | index (12 bit)
每毫秒最多生成4096个id，集群机器最多1024台
*/

type Snowflake struct {
	lastTimestamp int64
	index         int16
	machId        int16
}

func NewSnowflake(id int16) *Snowflake {
	sf := &Snowflake{}
	sf.Init(id)
	return sf
}

func (s *Snowflake) Init(id int16) error {
	if id > 0xff {
		return errors.New("illegal machine id")
	}

	s.machId = id
	s.lastTimestamp = time.Now().UnixNano() / 1e6
	s.index = 0
	return nil
}

func (s *Snowflake) GetId() (int64, error) {
	curTimestamp := time.Now().UnixNano() / 1e6
	if curTimestamp == s.lastTimestamp {
		s.index++
		if s.index > 0xfff {
			s.index = 0xfff
			return -1, errors.New("out of range")
		}
	} else {
		//fmt.Printf("id/ms:%d -- %d\n", s.lastTimestamp, s.index)
		s.index = 0
		s.lastTimestamp = curTimestamp
	}
	return int64((0x1ffffffffff&s.lastTimestamp)<<22) + int64(0xff<<10) + int64(0xfff&s.index), nil
}

var UUID *Snowflake
