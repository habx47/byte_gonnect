package router

import (
	"bytes"
	"io"
	"net"
	"testing"

	"github.com/habx47/byte_gonnect/parser"
)

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
}
