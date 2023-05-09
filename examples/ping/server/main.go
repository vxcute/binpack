package main

import (
	"flag"
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
	Opcode 	uint16
	Data    string
}

func handle(conn net.Conn) {

	defer conn.Close()

	buf := make([]byte, 2048)

	for {

		n, err := conn.Read(buf)

		if err != nil {
			break
		}

		var msg Message 

		if err := binpack.Unpack(buf[:n], &msg); err != nil {
			log.Println("err: ", err)
		}

		if err != nil { 
			log.Println("err: ", err)
			break
		}

		switch msg.Opcode {
		case PING:
			log.Printf("recevied ping message from: %s - MSG: %s\n", conn.RemoteAddr().String(), msg.Data)
				
			resp := Message{ 
				Opcode: PONG,
				Data: "PONG",
			}

			 buf, err := binpack.Pack(resp);

			 if err != nil {
				log.Println("Err: ", err)
				break
			}

			conn.Write(buf)
		}
	}
}

func main() {

	flag.Parse()

	listener, err := net.Listen("tcp", *addr)

	if err != nil {
		log.Fatal(err)
	} 

	log.Println("listening on: ", listener.Addr().String())

	for {
		conn, err := listener.Accept() 

		if err != nil {
			log.Println("err: ", err)
			continue
		}

		go handle(conn)
	}
}