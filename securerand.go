package main

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
)

func newRand() *mathrand.Rand {
	return mathrand.New(&ReaderSource{Reader: cryptorand.Reader})
}
