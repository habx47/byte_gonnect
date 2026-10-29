package router

import (
	"errors"
	"fmt"
	"net"

	"github.com/habx47/byte_gonnect/parser"
)

var ConnWriteErr = errors.New("failed to write to conn")

var GETresponseStr_root = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 25\r\nConnection: close\r\n\r\nHello, from the root page")
var GETresponseStr_hello = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 26\r\nConnection: close\r\n\r\nHello, from the hello page")
var GETresponseStr_default = []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\nContent-Length: 14 \r\nConnection: close\r\n\r\nPage not found")

var POSTresponseStr = []byte("HTTP/1.1 405 Method Not Allowed\r\nContent-Type: text/plain\r\nContent-Length: 18\r\nConnection: close\r\n\r\nMethod not allowed")

func Router(conn net.Conn, reqMetadata parser.ReqMetadata) error {
	method := reqMetadata.Method

	switch method {
	case "GET":
		err := handleGETReq(conn, reqMetadata)
		if err != nil {
			return err
		}
	case "POST":
		err := handlePOSTReq(conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleGETReq(conn net.Conn, reqMetadata parser.ReqMetadata) error {
	path := reqMetadata.Path

	switch path {
	case "/":
		_, err := conn.Write(GETresponseStr_root)
		if err != nil {
			return fmt.Errorf("%w: %v", ConnWriteErr, err)
		}
	case "/hello":
		_, err := conn.Write(GETresponseStr_hello)
		if err != nil {
			return fmt.Errorf("%w: %v", ConnWriteErr, err)
		}
	default:
		_, err := conn.Write(GETresponseStr_default)
		if err != nil {
			return fmt.Errorf("%w: %v", ConnWriteErr, err)
		}
	}

	return nil
}

func handlePOSTReq(conn net.Conn) error {
	_, err := conn.Write(POSTresponseStr)
	if err != nil {
		return fmt.Errorf("%w: %v", ConnWriteErr, err)
	}

	return nil
}
