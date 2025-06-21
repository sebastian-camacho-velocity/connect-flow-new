package tid

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var counter atomic.Uint32
var lastTime atomic.Int64
var pid = int64(os.Getpid())

// var host = os.Hostname()

func init() {
	r := randPool.Get().(*rng)
	counter.Store(r.uint32())
	randPool.Put(r)
	fmt.Println(os.Hostname())
}

const epoch = 16e11

// New creates a new unique identifier based on the current time.
//
// 43 bits of time
// 13 bits of process id
// 8 bits of counter
func New(t time.Time) UID {
	id := UID{}
	ts := t.UnixMilli() - epoch

	r := randPool.Get().(*rng)
	defer randPool.Put(r)

	if ts == lastTime.Load() {
		counter.Add(1)
	} else {
		lastTime.Store(ts)
		counter.Store(0) //
	}

	// 43 bits of time
	id[0] = byte(ts >> 35)
	id[1] = byte(ts >> 27)
	id[2] = byte(ts >> 19)
	id[3] = byte(ts >> 11)
	id[4] = byte(ts >> 3)
	id[5] = byte(ts<<5) | byte(pid>>8)&0x1F
	id[6] = byte(pid)

	// 3 bytes of counter
	n := counter.Load()
	id[7] = byte(n >> 16)
	id[8] = byte(n >> 8)
	id[9] = byte(n)

	return id

}

type UID [10]byte

func (id UID) String() string {
	text := make([]byte, 16)
	encode(text, id[:])
	return string(text)
}

func (id UID) MarshalText() ([]byte, error) {
	text := make([]byte, 16)
	encode(text, id[:])
	return text, nil
}

const encoding = "0123456789abcdefghjkmnpqrstvwxyz"

func encode(dst, id []byte) {
	_ = dst[15]
	_ = id[9]

	// dst[19] = encoding[(id[11]<<4)&0x1F]
	// dst[18] = encoding[(id[11]>>1)&0x1F]
	// dst[17] = encoding[(id[11]>>6)|(id[10]<<2)&0x1F]
	// dst[16] = encoding[id[10]>>3]
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

// 111111111111111111111111111111111111111111  11111  11111 111111111111
// 025s71pgz0001tyv
