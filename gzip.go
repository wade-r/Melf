package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net"
	"os"
)

func handleGzip(buf []byte, out *net.UDPConn) error {
	r, err := gzip.NewReader(bytes.NewBuffer(buf))
	defer r.Close()

	//DEBUG: output to stdout
	if err == nil {
		io.Copy(os.Stdout, r)
		r.Close()
	}

	//DEBUG: just relay to target
	out.Write(buf)

	return nil
}
