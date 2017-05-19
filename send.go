package main

import (
	"net"
)

var PacketMaxSize 	= 8 * 1024 // max bytes in a single UDP packet
var ChunckHeadSize	= 12

func sendGELF(buf []byte, out *net.UDPConn) error {
	if len(buf) < PacketMaxSize {
		out.Write(buf)
	} else {
	}
	return nil
}
