package router

import (
	"errors"
	"fmt"
	"net"

	"github.com/habx47/byte_gonnect/parser"
)

var ConnWriteErr = errors.New("failed to write to conn")
var MethodHandlingErr = errors.New("missing valid method handling")

type Routes map[string]bool
type EndpointResps map[string][]byte

var Rsp_NotFound = []byte("HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\nContent-Length: 14\r\nConnection: close\r\n\r\nPage not found")
var Rsp_MethodNotAllowed = []byte("HTTP/1.1 405 Method Not Allowed\r\nContent-Type: text/plain\r\nContent-Length: 18\r\nConnection: close\r\n\r\nMethod not allowed")

func Router(conn net.Conn, reqMetadata parser.ReqMetadata, routes Routes, endpointResps EndpointResps) error {

	if _, ok := routes[reqMetadata.Path]; !ok {
		_, err := conn.Write(Rsp_NotFound)
		if err != nil {
			return fmt.Errorf("%w: %v", ConnWriteErr, err)
		}
		return nil
	}

	endpoint := reqMetadata.Method + reqMetadata.Path
	if _, ok := endpointResps[endpoint]; !ok {
		_, err := conn.Write(Rsp_MethodNotAllowed)
		if err != nil {
			return fmt.Errorf("%w: %v", ConnWriteErr, err)
		}
		return nil
	}

	switch reqMetadata.Method {
	case "GET":
		err := handleGETReq(conn, reqMetadata, endpointResps)
		if err != nil {
			return err
		}
	case "POST":
		err := handlePOSTReq(conn, reqMetadata, endpointResps)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%v", MethodHandlingErr)
	}
	return nil
}

func handleGETReq(conn net.Conn, reqMetadata parser.ReqMetadata, endpointResps EndpointResps) error {
	endpoint := reqMetadata.Method + reqMetadata.Path

	_, err := conn.Write(endpointResps[endpoint])
	if err != nil {
		return fmt.Errorf("%w: %v", ConnWriteErr, err)
	}
	return nil
}

func handlePOSTReq(conn net.Conn, reqMetadata parser.ReqMetadata, endpointResps EndpointResps) error {
	endpoint := reqMetadata.Method + reqMetadata.Path

	_, err := conn.Write(endpointResps[endpoint])
	if err != nil {
		return fmt.Errorf("%w: %v", ConnWriteErr, err)
	}
	return nil
}
