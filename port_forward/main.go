package main

import (
	"io"
	"log"
	"net"
)

func forward(conn net.Conn) {
	client, err := net.Dial("tcp", "")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	go func() {
		_, err := io.Copy(client, conn)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	_, err = io.Copy(conn, client)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":40001")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go forward(conn)
	}
}
