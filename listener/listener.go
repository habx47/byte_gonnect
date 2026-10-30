package listener

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"github.com/habx47/byte_gonnect/endpoints"
	"github.com/habx47/byte_gonnect/parser"
	"github.com/habx47/byte_gonnect/router"
)

var URIErr = errors.New("URI too long")
var HeaderErr = errors.New("headers too large")

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
	routes := make(endpoints.Routes)
	endpointResps := endpoints.EndpointResps{
		"GET/":      endpoints.GETRsp_root,
		"GET/hello": endpoints.GETRsp_hello,
	}

	err := createRoutes(endpointResps, routes)
	if err != nil {
		logger.Fatal(err)
	}

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

		go handleConn(conn, logger, routes, endpointResps)
	}
}

func handleConn(conn net.Conn, logger Logger, routes endpoints.Routes, endpointResps endpoints.EndpointResps) {
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
			if isPrefix {
				_, err := conn.Write([]byte("HTTP/1.1 414 URI too long\r\nConnection: close\r\n\r\n"))
				if err != nil {
					logger.Printf("failed to write to conn: %v", err)
				}
				logger.Printf("%v", URIErr)
				return
			}

			reqLineErr := reqMetadata.ParseRequestLine(line)
			if reqLineErr != nil {
				_, err := conn.Write(endpoints.Rsp_BadRequest)
				if err != nil {
					logger.Printf("failed to write to client: %v", err)
				}
				logger.Printf("failed to parse the request: %v", reqLineErr)
				return
			}
			state = stateHeaders
		case stateHeaders:
			if isPrefix {
				_, err := conn.Write([]byte("HTTP/1.1 431 Request Header Fields Too Large\r\nConnection: close\r\n\r\n"))
				if err != nil {
					logger.Printf("failed to write to conn: %v", err)
				}
				logger.Printf("%v", HeaderErr)
				return
			}
			parseHeadersErr := reqMetadata.ParseHeaders(line)
			if parseHeadersErr != nil {
				conn.Write(endpoints.Rsp_BadRequest)
				logger.Printf("failed to parse request: %v", parseHeadersErr)
				return
			}
		}
	}

	if reqMetadata.ContentLength.Value > MaxRequestBodySize {
		_, err := conn.Write([]byte("HTTP/1.1 413 Payload Too Large\r\nConnection: close\r\n\r\n"))
		if err != nil {
			logger.Printf("failed to write to conn: %v", err)
		}
		logger.Printf("%v: request body size more than allowed", parser.BadRequestErr)
		return
	}

	body := make([]byte, reqMetadata.ContentLength.Value)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			_, err := conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nConnection: close\r\n\r\n"))
			if err != nil {
				logger.Printf("failed to write to conn: %v", err)
			}
		} else {
			logger.Printf("failed to read request body: %v", err)
		}
		return
	}

	routerErr := router.Router(conn, reqMetadata, routes, endpointResps)
	if routerErr != nil {
		logger.Printf("failed to write to client: %v", routerErr)
		return
	}
}
