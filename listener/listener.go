package listener

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"github.com/habx47/byte_gonnect/parser"
	"github.com/habx47/byte_gonnect/router"
)

type Logger interface {
	Printf(format string, v ...any)
}

const (
	stateRequestLine int = iota
	stateHeaders
)

const MaxRequestBodySize = 10000

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

	reader := bufio.NewReaderSize(conn, 8192)
	for {
		line, isPrefix, connReadErr := reader.ReadLine()
		if connReadErr != nil {
			logger.Printf("client disconnected or error: %v", connReadErr)
			return
		}

		if state == stateHeaders && len(line) == 0 {
			break
		}

		switch state {
		case stateRequestLine:
			reqLineErr := reqMetadata.ParseRequestLine(line, isPrefix)
			if errors.Is(reqLineErr, parser.URIErr) {
				conn.Write([]byte("HTTP/1.1 414 URI too long\r\n\r\n"))
				logger.Printf("%v", parser.URIErr)
				return
			}
			if reqLineErr != nil {
				logger.Printf("failed to parse request: %v", reqLineErr)
				return
			}
			state = stateHeaders
		case stateHeaders:
			parseHeadersErr := reqMetadata.ParseHeaders(line, isPrefix)
			if errors.Is(parseHeadersErr, parser.HeaderErr) {
				conn.Write([]byte("HTTP/1.1 431 Request Header Fields Too Large\r\n\r\n"))
				logger.Printf("%v", parser.HeaderErr)
				return
			}
			if parseHeadersErr != nil {
				logger.Printf("failed to parse request: %v", parseHeadersErr)
				return
			}
		}
	}

	if reqMetadata.ContentLength.Value > MaxRequestBodySize {
		conn.Write([]byte("HTTP/1.1 413 Payload Too Large\r\nConnection: close\r\n\r\n"))
		logger.Printf("%v: request body size more than allowed", parser.MalformedReqErr)
		return
	}

	body := make([]byte, reqMetadata.ContentLength.Value)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"))
		} else {
			logger.Printf("failed to read request body: %v", err)
		}
		return
	}

	connWriteErr := router.Router(conn, reqMetadata)
	if connWriteErr != nil {
		logger.Printf("failed to write to client: %v", connWriteErr)
		return
	}
}
