package uid

import (
	"sync/atomic"
	"time"
	"unsafe"
)

var counter atomic.Uint32

func init() {
	r := randPool.Get().(*rng)
	counter.Store(r.uint32())
	randPool.Put(r)
}

func New() string {
	return NewWithTime(time.Now().Unix() - 15e8)
}

func NewWithTime(ts int64) string {
	var id [10]byte

	// 4 bytes of time
	id[0] = byte(ts >> 24)
	id[1] = byte(ts >> 16)
	id[2] = byte(ts >> 8)
	id[3] = byte(ts)

	// 3 bytes of counter
	n := counter.Add(1)
	id[4] = byte(n >> 16)
	id[5] = byte(n >> 8)
	id[6] = byte(n)

	// 3 bytes of random
	r := randPool.Get().(*rng)
	n = r.uint32()
	id[7] = byte(n >> 16)
	id[8] = byte(n >> 8)
	id[9] = byte(n)
	randPool.Put(r)

	text := make([]byte, 16)
	encode(text, id[:])

	return *(*string)(unsafe.Pointer(&text))
}

const encoding = "0123456789abcdefghjkmnpqrstvwxyz"

func encode(dst, id []byte) {
	_ = dst[15]
	_ = id[9]

	dst[15] = encoding[id[9]&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[11] = encoding[(id[7]>>4)|(id[6]<<4)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[9] = encoding[(id[6]>>6)|(id[5]<<2)&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[7] = encoding[id[4]&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[3] = encoding[(id[2]>>4)|(id[1]<<4)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[1] = encoding[(id[1]>>6)|(id[0]<<2)&0x1F]
	dst[0] = encoding[id[0]>>3]
}
