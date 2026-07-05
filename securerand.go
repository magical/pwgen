package main

import (
	"crypto/rand"
	"math/bits"
)

func randUint64() uint64 {
	var buf [64 / 8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		panic("readUint64: " + err.Error())
	}
	entropy.Add(64)
	v := uint64(buf[0])<<0 + uint64(buf[1])<<8 + uint64(buf[2])<<16 + uint64(buf[3])<<24 +
		uint64(buf[4])<<32 + uint64(buf[5])<<40 + uint64(buf[6])<<48 + uint64(buf[7])<<56
	return v
}

func randInt(n int) int {
	if n < 0 {
		panic("n out of range")
	}
	if n&(n-1) == 0 {
		return int(randUint64()) & (n - 1)
	}
	// a fast, efficient, and completely unbiased method for
	// sampling a uniform value, as described by Daniel Lemire.
	u := uint64(n)
	hi, lo := bits.Mul64(randUint64(), u)
	if lo < u {
		thresh := -u % u
		for lo < thresh {
			hi, lo = bits.Mul64(randUint64(), u)
		}
	}
	return int(hi)
}

type SecureRand struct{}

func (*SecureRand) Intn(n int) int { return randInt(n) }

func newRand() *SecureRand {
	return new(SecureRand)
}
