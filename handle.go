package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"os"
)

func handleGzip(conn net.PacketConn, sendAddr net.Addr, buf []byte, recvAddr net.Addr) error {
	// create a reader
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	defer r.Close()

	//DEBUG: output to stdout
	if err == nil {
		io.Copy(os.Stdout, r)
	}

	// extract short_message

	//DEBUG: just relay to target
	conn.WriteTo(buf, sendAddr)

	return nil
}

// for format other than gzip, just relay
func handleUnknown(conn net.PacketConn, sendAddr net.Addr, buf []byte, recvAddr net.Addr) {
	// make a local copy, since 'buf' is a subslice of shared slice
	local := make([]byte, len(buf))
	copy(buf, local)
	go conn.WriteTo(local, sendAddr)
}
