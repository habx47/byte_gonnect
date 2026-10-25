package listener

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"testing/iotest"
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
		inputMessage := "hello, world!\n"
		serverSide, clientSide := net.Pipe()
		defer clientSide.Close()

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		go handleConn(serverSide, logger)

		_, err := clientSide.Write([]byte(inputMessage))
		if err != nil {
			t.Fatalf("failed to write to server: %v", err)
		}

		reader := bufio.NewReader(clientSide)
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("failed to read response from server: %v", err)
		}

		expectedResponse := inputMessage
		if response != expectedResponse {
			t.Errorf("got %v, want %v", response, expectedResponse)
		}

		clientSide.Close()
	})

	t.Run("error on conn read", func(t *testing.T) {
		serverSide, _ := net.Pipe()
		customErr := errors.New("read err")
		faultyConn := &ReadErrConn{
			Conn:      serverSide,
			errReader: iotest.ErrReader(customErr),
		}

		var logBuffer bytes.Buffer
		logger := log.New(&logBuffer, "", 0)

		handleConn(faultyConn, logger)

		got := logBuffer.String()

		if !strings.Contains(got, "client disconnected or error") {
			t.Errorf("expected connReadErr branch to execute, but got: \n%s", got)
		}

		if !strings.Contains(got, customErr.Error()) {
			t.Errorf("expected log to contain %q, but got: \n%s", customErr.Error(), got)
		}
	})

	t.Run("error on conn write", func(t *testing.T) {
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
			handleConn(faultyConn, logger)
			close(done)
		}()

		_, err := clientSide.Write([]byte("hello, world!\n"))
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
}
