package listener

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/habx47/byte_gonnect/endpoints"
	"github.com/habx47/byte_gonnect/parser"
)

type ReadErrConn struct {
	net.Conn
	errReader io.Reader
}

func (r *ReadErrConn) Read(b []byte) (int, error) {
	return r.errReader.Read(b)
}

var ServerWriteErr = "server write error"

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

func TestHandleConn(t *testing.T) {
	t.Run("success on read & write", func(t *testing.T) {
		request := "GET / HTTP/1.1\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, Server"
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/":      endpoints.GETRsp_root,
			"GET/hello": endpoints.GETRsp_hello,
		}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		go handleConn(serverSide, logger, stubRoutes, stubEndpointResps)

		_, err := clientSide.Write([]byte(request))
		if err != nil {
			t.Fatalf("failed to write to server: %v", err)
		}

		response, err := io.ReadAll(clientSide)
		if err != nil {
			t.Fatalf("failed to read response from server: %v", err)
		}

		if !bytes.Equal(response, endpoints.GETRsp_root) {
			t.Errorf("got: %s,\n want: %s", response, endpoints.GETRsp_root)
		}

		clientSide.Close()
	})

	t.Run("error on conn read", func(t *testing.T) {
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/":      endpoints.GETRsp_root,
			"GET/hello": endpoints.GETRsp_hello,
		}
		serverSide, _ := net.Pipe()
		customErr := errors.New("read err")
		faultyConn := &ReadErrConn{
			Conn:      serverSide,
			errReader: iotest.ErrReader(customErr),
		}

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		handleConn(faultyConn, logger, stubRoutes, stubEndpointResps)

		got := logBuffer.String()

		if !strings.Contains(got, "client disconnected or error") {
			t.Errorf("expected connReadErr branch to execute, but got: \n%s", got)
		}

		if !strings.Contains(got, customErr.Error()) {
			t.Errorf("expected log to contain %q, but got: \n%s", customErr.Error(), got)
		}
	})

	t.Run("error on conn write", func(t *testing.T) {
		request := "GET / HTTP/1.1\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, Server"
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/":      endpoints.GETRsp_root,
			"GET/hello": endpoints.GETRsp_hello,
		}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()
		faultyConn := &WriteErrConn{
			Conn:        serverSide,
			faultyWrite: false,
		}

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		done := make(chan struct{})

		go func() {
			faultyConn.faultyWrite = true
			handleConn(faultyConn, logger, stubRoutes, stubEndpointResps)
			close(done)
		}()

		_, err := clientSide.Write([]byte(request))
		if err != nil {
			t.Fatalf("failed to write to server: %v", err)
		}

		<-done

		got := logBuffer.String()

		if !strings.Contains(got, "failed to write to client") {
			t.Errorf("expected connWriteErr branch to execute, but got: \n%s", got)
		}

		if !strings.Contains(got, ServerWriteErr) {
			t.Errorf("expected log to contain %q, but got: \n%s", ServerWriteErr, got)
		}
	})

	t.Run("malformed request line", func(t *testing.T) {
		request := "/ HTTP/1.1\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello, Server"
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/":      endpoints.GETRsp_root,
			"GET/hello": endpoints.GETRsp_hello,
		}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		done := make(chan struct{})

		go func() {
			handleConn(serverSide, logger, stubRoutes, stubEndpointResps)
			close(done)
		}()

		_, err := clientSide.Write([]byte(request))
		if err != nil {
			t.Fatalf("failed to write to server: %v", err)
		}

		<-done

		got := logBuffer.String()

		if !strings.Contains(got, "failed to parse request") {
			t.Errorf("expected reqLineErr branch to execute, but got :\n%s", got)
		}

		if !strings.Contains(got, parser.MalformedReqErr.Error()) {
			t.Errorf("expected log to contain %q, but got :\n%s", parser.MalformedReqErr, got)
		}

		clientSide.Close()
	})

	t.Run("malformed Content-Length", func(t *testing.T) {
		request := "GET / HTTP/1.1\r\nContent-Type: text/plain\r\nContent-Length: abc\r\n\r\nHello, Server"
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/":      endpoints.GETRsp_root,
			"GET/hello": endpoints.GETRsp_hello,
		}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		done := make(chan struct{})

		go func() {
			handleConn(serverSide, logger, stubRoutes, stubEndpointResps)
			close(done)
		}()

		_, err := clientSide.Write([]byte(request))
		if err != nil {
			t.Fatalf("failed to write to server: %v", err)
		}

		<-done

		got := logBuffer.String()

		if !strings.Contains(got, "failed to parse request") {
			t.Errorf("expected parseHeadersErr branch to execute, but got :\n%s", got)
		}

		if !strings.Contains(got, parser.ContentLengthErr.Error()) {
			t.Errorf("expected log to contain %q, but got :\n%s", parser.ContentLengthErr, got)
		}

		clientSide.Close()
	})

	t.Run("request line too long", func(t *testing.T) {
		r := []byte("this is a huge request")
		request := bytes.Repeat(r, (8192/len(r))+1)
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/": endpoints.GETRsp_root,
		}
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		go handleConn(serverSide, logger, stubRoutes, stubEndpointResps)

		go func() {
			clientSide.Write(request)
		}()

		response, err := io.ReadAll(clientSide)
		if err != nil {
			t.Fatalf("failed to read response from server: %v", err)
		}

		if !bytes.Equal(response, []byte("HTTP/1.1 414 URI too long\r\n\r\n")) {
			t.Errorf("got: %s,\n want: %v", response, "HTTP/1.1 414 URI too long\r\n\r\n")
		}

		clientSide.Close()
	})

}

func TestCreateRoutes(t *testing.T) {
	stubRoutes := map[string]bool{}
	stubEndpointResps := map[string][]byte{
		"GET/":      endpoints.GETRsp_root,
		"GET/hello": endpoints.GETRsp_hello,
	}

	tests := []struct{ expected string }{
		{expected: "/"},
		{expected: "/hello"},
	}

	err := createRoutes(stubEndpointResps, stubRoutes)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if len(stubRoutes) != len(stubEndpointResps) {
		t.Fatalf("mismatch between routes registered and created. got: %v, registered: %v", len(stubRoutes), len(stubEndpointResps))
	}

	for _, tt := range tests {
		exists := stubRoutes[tt.expected]
		if !exists {
			t.Errorf("expected route: %v not found", tt.expected)
		}
	}
}
