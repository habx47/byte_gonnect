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
}
