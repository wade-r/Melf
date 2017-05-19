package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"github.com/tidwall/gjson"
	"log"
	"net"
)

// decompress, handle it
// WARN: buf is a copy, do anything you want
// WANR: this func will run in a go routine
func handleGzip(sendAddr net.Addr, buf []byte, recvAddr net.Addr) {
	// create a gzip.Reader
	gr, err := gzip.NewReader(bytes.NewBuffer(buf))
	defer gr.Close()

	if err != nil {
		log.Println("failed to create gzip reader")
		return
	}

	// create a bufio.Reader over gzip.Reader
	r := bufio.NewReader(gr)

	// read the string
	s, _ := r.ReadString(0)

	// parse json
	m, ok := gjson.Parse(s).Value().(map[string]interface{})

	if !ok {
		log.Println("failed to parse GELF json:", s)
		return
	}

	// find cid
	cid, ok := m[kContainerId].(string)

	if len(cid) == 0 {
		// container_id not found, just relay
		_, err := sendPacket(buf)
		if err != nil {
			log.Println("failed to relay gzip message")
			log.Println(err)
		}
	} else {
		// dispatch message to sessions
		dispatchMessageToSessions(m)
	}
}

// for format other than gzip, just relay
// WARN: buf is a copy, do anything you want
// WANR: this func will run in a go routine
func handleWild(sendAddr net.Addr, buf []byte, recvAddr net.Addr) {
	_, err := sendPacket(buf)
	if err != nil {
		log.Println("failed to relay wild message")
		log.Println(err)
	}
}
