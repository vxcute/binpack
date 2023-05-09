package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/vxcute/binpack"
)

var ( 
	addr = flag.String("addr", ":1337", "server address")
)

const ( 
	PING = 0x10
	PONG = 0x12
)

type Message struct {
	Opcode uint16	
	Data   string 
}

func main() {

	flag.Parse()

	conn, err := net.Dial("tcp", *addr) 

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	msg := Message{ 
		Opcode: PING,
		Data: "PING",
	}


	err = binpack.Pack(conn, msg)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("wrote message")

	if  err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 2048)

	n, err := conn.Read(buf)
	
	if err != nil {
		log.Fatal(err)
	}

	var m Message 

	if err := binpack.Unpack(buf[:n], &m); err != nil {
		log.Fatal(err)
	}

	fmt.Println(m)
}
