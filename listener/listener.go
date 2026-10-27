package listener

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os"

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
			parser.ParseRequestLine(line, reqMetadata)
			state = stateHeaders
		case stateHeaders:
			parser.ParseHeaders(line, reqMetadata)
		}
	}

	/*
	   _, connWriteErr := conn.Write(line)

	   	if connWriteErr != nil {
	   		logger.Printf("failed to write to client: %v", connWriteErr)
	   		return
	   	}
	*/
}
