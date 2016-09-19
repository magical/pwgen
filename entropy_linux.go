package main

import (
	"expvar"
	"io/ioutil"
)

func init() {
	expvar.Publish("entropy_avail", expvar.Func(func() interface{} {
		return getEntropyEstimate()
	}))
}

func getEntropyEstimate() int64 {
	b, err := ioutil.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		return -1
	}
	return atoi(b)
}

func atoi(b []byte) int64 {
	var v int64
	for i, c := range b {
		if '0' <= c && c <= '9' {
			v = v*10 + int64(c) - '0'
			if v < 0 {
				return -1
			}
		} else if c == '\n' && i == len(b)-1 {
			// nothing
		} else {
			return -1
		}
	}
	return v
}
