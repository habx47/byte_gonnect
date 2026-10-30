package router

import (
	"errors"
	"fmt"
	"net"

	"github.com/habx47/byte_gonnect/endpoints"
	"github.com/habx47/byte_gonnect/parser"
)

var ConnWriteErr = errors.New("failed to write to conn")
var MethodHandlingErr = errors.New("registered method handling not implemented")

func Router(conn net.Conn, reqMetadata parser.ReqMetadata, routes endpoints.Routes, endpointResps endpoints.EndpointResps) error {

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

func handleGETReq(conn net.Conn, reqMetadata parser.ReqMetadata, endpointResps endpoints.EndpointResps) error {
	endpoint := reqMetadata.Method + reqMetadata.Path

	_, err := conn.Write(endpointResps[endpoint])
	if err != nil {
		return fmt.Errorf("%w: %v", ConnWriteErr, err)
	}
	return nil
}

func handlePOSTReq(conn net.Conn, reqMetadata parser.ReqMetadata, endpointResps endpoints.EndpointResps) error {
	endpoint := reqMetadata.Method + reqMetadata.Path

	_, err := conn.Write(endpointResps[endpoint])
	if err != nil {
		return fmt.Errorf("%w: %v", ConnWriteErr, err)
	}
	return nil
}
