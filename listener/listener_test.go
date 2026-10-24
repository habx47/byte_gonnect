package listener

import (
	"bufio"
	"net"
	"testing"
)

func TestHandleConn(t *testing.T) {
	inputMessage := "hello, world!\n"
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	go handleConn(serverSide)

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
}
