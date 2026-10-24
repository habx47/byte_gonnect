package listener

import (
	"bufio"
	"log"
	"net"
)

func Listener() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept a connection: %v", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("client disconnected on err: %v", err)
			return
		}

		_, connWriteErr := conn.Write(message)
		if connWriteErr != nil {
			log.Printf("failed to write to client: %v", connWriteErr)
			return
		}
	}
}
