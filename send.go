package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
)

var PacketMaxSize = 8 * 1024 // max bytes in a single UDP packet
var ChunckHeadSize = 12

// WARN: m is not a shared object
// WARN: will run in go routine
func sendMessage(m Message) {
	// marshal json
	raw, err := json.Marshal(m)
	if err != nil {
		log.Println("failed to marshal:", m)
		return
	}
	// prepare buffer
	buf := &bytes.Buffer{}
	// write gzip
	gr := gzip.NewWriter(buf)
	gr.Write(raw)
	gr.Close()
	// send
	out := buf.Bytes()
	if len(out) < PacketMaxSize {
		sendPacket(out)
	} else {
		//TODO: implement chuncked format
	}
}
