package main

import (
	xhmac "crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
)

func hmac(key, data []byte) []byte {
	h := xhmac.New(sha1.New, key)
	h.Write(data)
	b := make([]byte, 0, sha1.Size)
	return h.Sum(b)
}

func hotp(key []byte, counter uint64, digits int) int {
	ctr := make([]byte, 8)
	binary.BigEndian.PutUint64(ctr, counter)
	sum := hmac(key, ctr)
	idx := sum[len(sum)-1] & 0x0F
	v := binary.BigEndian.Uint32(sum[idx:]) & 0x7FFFFFFF
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}
	return int(v % d)
}
