package uid

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"
)

// RandString generates a random string of length n.
func RandString(charset string, n uint8) string {
	l := uint32(len(charset))
	b := make([]byte, n)
	r := randPool.Get().(*rng)
	rn := r.uint32()
	for i := range b {
		if rn == 0 {
			rn = r.uint32()
		}
		b[i] = charset[rn%l]
		rn = rn / l
	}
	randPool.Put(r)
	return string(b)
}

var randPool = sync.Pool{
	New: func() interface{} {
		return &rng{
			n: genRandSeed(),
		}
	},
}

func genRandSeed() uint32 {
	b := make([]byte, 8)
	rand.Read(b)
	x := binary.BigEndian.Uint64(b) + uint64(time.Now().UnixNano())
	return uint32((x >> 32) ^ x)
}

type rng struct {
	n uint32
}

func (r *rng) uint32() uint32 {
	// Xorshift https://en.wikipedia.org/wiki/Xorshift
	n := r.n
	n ^= n << 13
	n ^= n >> 17
	n ^= n << 5
	r.n = n
	return n
}

const (
	CharsetAlphaNum   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CharsetAlpha      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetAlphaLower = "abcdefghijklmnopqrstuvwxyz"
	CharsetAlphaUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetNum        = "0123456789"
	CharsetHex        = "0123456789abcdef"
	CharsetHexUpper   = "0123456789ABCDEF"
)
