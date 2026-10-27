package listener

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/habx47/byte_gonnect/parser"
)

type Logger interface {
	Printf(format string, v ...any)
}

const (
	stateRequestLine int = iota
	stateHeaders
)

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
	state := stateRequestLine
	responseStr := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, Client"
	reqMetadata := parser.ReqMetadata{}

	reader := bufio.NewReader(conn)
	for {
		line, connReadErr := reader.ReadBytes('\n')
		if connReadErr != nil {
			logger.Printf("client disconnected or error: %v", connReadErr)
			return
		}

		line = bytes.TrimRight(line, "\r\n")

		if state == stateHeaders && len(line) == 0 {
			break
		}

		switch state {
		case stateRequestLine:
			reqMetadata.ParseRequestLine(line)
			state = stateHeaders
		case stateHeaders:
			reqMetadata.ParseHeaders(line)
		}
	}

	if reqMetadata.ContentLength > 0 {
		body := make([]byte, reqMetadata.ContentLength)
		if _, err := io.ReadFull(reader, body); err != nil {
			logger.Printf("failed to read request body: %v", err)
			return
		}
	}

	_, connWriteErr := conn.Write([]byte(responseStr))

	if connWriteErr != nil {
		logger.Printf("failed to write to client: %v", connWriteErr)
		return
	}
}
