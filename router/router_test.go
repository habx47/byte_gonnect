package router

import (
	"bytes"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/habx47/byte_gonnect/endpoints"
	"github.com/habx47/byte_gonnect/parser"
)

const ServerWriteErr = "server write error"

var Rsp_root = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 25\r\nConnection: close\r\n\r\nHello, from the root page")

type WriteErrConn struct {
	net.Conn
	faultyWrite bool
}

func (w *WriteErrConn) Write(b []byte) (n int, err error) {
	if w.faultyWrite {
		return 0, errors.New(ServerWriteErr)
	}
	return w.Conn.Write(b)
}

func TestGETReqHandling(t *testing.T) {
	t.Run("success on write to conn on path: /", func(t *testing.T) {
		stubReqMetadata := parser.ReqMetadata{Method: "GET", Path: "/"}
		stubEndpointResps := endpoints.EndpointResps{"GET/": Rsp_root}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		go func() {
			handleGETReq(serverSide, stubReqMetadata, stubEndpointResps)
			serverSide.Close()
		}()

		response, err := io.ReadAll(clientSide)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		if !bytes.Equal(response, Rsp_root) {
			t.Errorf("got: %s,\n want: %s", response, Rsp_root)
		}

		clientSide.Close()
	})

	t.Run("error on write", func(t *testing.T) {
		stubReqMetadata := parser.ReqMetadata{Method: "GET", Path: "/"}
		stubEndpointResps := endpoints.EndpointResps{"GET/": Rsp_root}
		serverSide, _ := net.Pipe()
		defer serverSide.Close()
		faultyConn := WriteErrConn{
			serverSide,
			true,
		}

		err := handleGETReq(&faultyConn, stubReqMetadata, stubEndpointResps)
		if err == nil {
			t.Fatalf("expected an error: %v, but got nil", ConnWriteErr)
		}
		if !errors.Is(err, ConnWriteErr) {
			t.Errorf("got: %v, want: %v", err, ConnWriteErr)
		}

		serverSide.Close()
	})
}

func TestPOSTReqHandling(t *testing.T) {
	t.Run("success on write to conn", func(t *testing.T) {
		stubReqMetadata := parser.ReqMetadata{Method: "POST", Path: "/"}
		stubEndpointResps := endpoints.EndpointResps{"POST/": Rsp_root}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		go func() {
			handlePOSTReq(serverSide, stubReqMetadata, stubEndpointResps)
			serverSide.Close()
		}()

		response, err := io.ReadAll(clientSide)
		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}

		if !bytes.Equal(response, Rsp_root) {
			t.Errorf("got: %s,\n want: %s", response, Rsp_root)
		}
	})

	t.Run("error on write", func(t *testing.T) {
		stubReqMetadata := parser.ReqMetadata{Method: "POST", Path: "/"}
		stubEndpointResps := endpoints.EndpointResps{"POST/": Rsp_root}
		serverSide, _ := net.Pipe()
		defer serverSide.Close()

		faultyConn := WriteErrConn{
			serverSide,
			true,
		}

		err := handlePOSTReq(&faultyConn, stubReqMetadata, stubEndpointResps)
		if err == nil {
			t.Fatalf("expected an error: %v, but got nil", ConnWriteErr)
		}
		if !errors.Is(err, ConnWriteErr) {
			t.Errorf("got: %v, want: %v", err, ConnWriteErr)
		}

		serverSide.Close()
	})
}
