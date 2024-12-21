package main

import (
	"log"
	"net"
	"sync"
)

var helloWorldResponse = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\nConnection: close\r\n\r\nHello, World!")
var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("Server is running on :3000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	n, err := conn.Read(buf)
	if err != nil || n < 4 || string(buf[:4]) != "GET " {
		conn.Write([]byte("HTTP/1.1 405 Method Not Allowed\r\n\r\n"))
		return
	}

	conn.Write(helloWorldResponse)
}
