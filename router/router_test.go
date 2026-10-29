package router

import (
	"bytes"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/habx47/byte_gonnect/parser"
)

const ServerWriteErr = "server write error"

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
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		go handleGETReq(serverSide, stubReqMetadata)

		response, err := io.ReadAll(clientSide)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		if !bytes.Equal(response, GETresponseStr_root) {
			t.Errorf("got: %s,\n want: %s", response, GETresponseStr_root)
		}

		clientSide.Close()
	})

	t.Run("error on write", func(t *testing.T) {
		stubReqMetadata := parser.ReqMetadata{Method: "GET", Path: "/"}
		serverSide, _ := net.Pipe()

		faultyConn := WriteErrConn{
			serverSide,
			true,
		}

		err := handleGETReq(&faultyConn, stubReqMetadata)

		if err == nil {
			t.Fatalf("expected an error: %v, but got nil", ConnWriteErr)
		}

		if !errors.Is(err, ConnWriteErr) {
			t.Errorf("got: %v, want: %v", err, ConnWriteErr)
		}
	})
}
