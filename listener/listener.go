package listener

import (
	"bufio"
	"log"
	"net"
	"os"
)

type Logger interface {
	Printf(format string, v ...any)
}

func Listener(port string) {
	logger := log.New(os.Stderr, "[SERVER] ", log.LstdFlags|log.Lshortfile)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Printf("failed to accept a connection: %v", err)
			continue
		}

		go handleConn(conn, logger)
	}
}

func handleConn(conn net.Conn, logger Logger) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, connReadErr := reader.ReadBytes('\n')
		if connReadErr != nil {
			logger.Printf("client disconnected or error: %v", connReadErr)
			return
		}

		_, connWriteErr := conn.Write(message)
		if connWriteErr != nil {
			logger.Printf("failed to write to client: %v", connWriteErr)
			return
		}
	}
}
