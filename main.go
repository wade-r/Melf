package main

import "flag"
import "net"
import "log"

// Constants
const kContainerId = "_container_id"
const kShortMessage = "short_message"

// Options
var bindOption = flag.String("b", "0.0.0.0:12201", "bind address")
var trgtOption = flag.String("t", "127.0.0.1:12202", "target address")

// Globals
var conn net.PacketConn
var sendAddr *net.UDPAddr

func main() {
	var err error

	flag.Parse()

	/*********************************
	* Resolve Target Address
	*********************************/

	sendAddr, err = net.ResolveUDPAddr("udp", *trgtOption)

	if err != nil {
		log.Println("invalid target address:", *trgtOption)
		log.Fatalln(err)
	}

	/*********************************
	* Create PacketConn
	*********************************/

	// use PacketConn for strict packet control and high performance
	conn, err = net.ListenPacket("udp", *bindOption)

	if err != nil {
		log.Println("cannot bind address:", *bindOption)
		log.Fatalln(err)
	}

	log.Println("melf running", *bindOption, "->", *trgtOption)

	/*********************************
	* Run Session Management Loop
	*********************************/

	go sessionManagementLoop()

	/*********************************
	* Run RecvLoop
	*********************************/

	// 16K is far enough for UDP packet receiving
	var buf = make([]byte, 2*PacketMaxSize)

	for {
		n, recvAddr, err := conn.ReadFrom(buf)

		if err != nil {
			log.Println("failed to read UDP:", err)
			continue
		}

		// Need first 2 byte for magic number testing
		if n < 2 {
			log.Println("UDP message is too short")
			continue
		}

		// make a local copy
		lbuf := make([]byte, n)
		copy(lbuf, buf[:n])

		// Detect GZIP
		// https://github.com/Graylog2/graylog2-server/blob/master/graylog2-server/src/main/java/org/graylog2/inputs/codecs/gelf/GELFMessage.java
		// RFC 1952
		if buf[0] == 0x1F && buf[1] == 0x8B {
			go handleGzip(sendAddr, lbuf, recvAddr)
		} else {
			go handleWild(sendAddr, lbuf, recvAddr)
		}
	}
}

func sendPacket(buf []byte) (int, error) {
	return conn.WriteTo(buf, sendAddr)
}
