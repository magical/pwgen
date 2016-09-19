package main

import (
	"flag"
	"net"
	"os"
	"strconv"
	"syscall"
)

var bind = flag.String("bind", ":8080", "address to listen on")

func listenfd() *os.File {
	pid, err := strconv.Atoi(os.Getenv("LISTEN_PID"))
	if err != nil || pid != os.Getpid() {
		return nil
	}
	nfds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || nfds != 1 {
		return nil
	}
	const fd = 3
	syscall.CloseOnExec(fd)
	return os.NewFile(uintptr(fd), "LISTEN_FD_3")
}

func listen() (net.Listener, error) {
	if f := listenfd(); f != nil {
		return net.FileListener(f)
	}

	return net.Listen("tcp", *bind)
}
