package main

import "flag"
import "net"
import "log"

var PacketMaxSize = 8 * 1024 // max bytes in a single UDP packet

var bindOption = flag.String("b", "0.0.0.0:12201", "bind address")
var trgtOption = flag.String("t", "127.0.0.1:12202", "target address")
var ruleOption = flag.String("r", "W", "rules")

func main() {
	flag.Parse()

	/*********************************
	* Resolve Addresses
	*********************************/

	recvAddr, err := net.ResolveUDPAddr("udp", *bindOption)

	if err != nil {
		log.Println("invalid bind address:", *bindOption)
		log.Fatalln(err)
	}

	sendAddr, err := net.ResolveUDPAddr("udp", *trgtOption)

	if err != nil {
		log.Println("invalid target address:", *trgtOption)
		log.Fatalln(err)
	}

	/*********************************
	* Create UDPConns
	*********************************/

	sendConn, err := net.DialUDP("udp", nil, sendAddr)

	if err != nil {
		log.Println("failed to dail target address:", *trgtOption)
		log.Fatalln(err)
	}

	recvConn, err := net.ListenUDP("udp", recvAddr)

	if err != nil {
		log.Println("failed to bind address:", *bindOption)
		log.Fatalln(err)
	}

	log.Println("melf running", *bindOption, "->", *trgtOption)

	/*********************************
	* Run RecvLoop
	*********************************/

	// 16K is far enough for UDP packet receiving
	var buf = make([]byte, 2*PacketMaxSize)

	for {
		n, addr, err := recvConn.ReadFromUDP(buf)

		if err != nil {
			log.Println("failed to read UDP:", err)
			continue
		}

		if n < 2 {
			log.Println("UDP message is too short")
			continue
		}

		// Detech CHUNCK
		if buf[0] == 0x1E && buf[1] == 0x0F {
			handleChunck(buf[:n], addr, sendConn)
		} else {
			handleGzip(buf[:n], sendConn)
		}
	}
}
