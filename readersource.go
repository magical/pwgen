package main

import (
	"io"
	"sync/atomic"
)

type ReaderSource struct {
	io.Reader
	buf [8]byte
	r   int64
}

func (s *ReaderSource) Seed(seed int64) {}
func (s *ReaderSource) Int63() int64 {
	_, err := io.ReadFull(s.Reader, s.buf[:8])
	if err != nil {
		panic("ReaderSource: read failed")
	}
	atomic.AddInt64(&s.r, 64)
	v := int64(s.buf[0])<<0 + int64(s.buf[1])<<8 + int64(s.buf[2])<<16 + int64(s.buf[3])<<24 +
		int64(s.buf[4])<<32 + int64(s.buf[5])<<40 + int64(s.buf[6])<<48 + int64(s.buf[7]&^0x80)<<56
	return v
}

// Bits returns the total number of bits read.
func (s *ReaderSource) Bits() int64 {
	return atomic.LoadInt64(&s.r)
}
