package main

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
)

var rand = mathrand.New(&ReaderSource{Reader: cryptorand.Reader})
