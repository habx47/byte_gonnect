package parser

import (
	"strings"
	"testing"
)

var StubMetadata ReqMetadata = make(ReqMetadata)

func TestRequestLine(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		tests := []struct {
			key string
		}{
			{"method"},
			{"path"},
			{"version"},
		}
		line := []byte("GET / HTTP/1.1")

		err := ParseRequestLine(line, StubMetadata)

		if err != nil {
			t.Fatalf("failed to parser request status line: %v", err)
		}

		for _, tt := range tests {
			if _, ok := StubMetadata[tt.key]; !ok {
				t.Fatalf("expected key: %v to exist", tt.key)
			}
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		line := []byte("GET /")

		err := ParseRequestLine(line, StubMetadata)

		if err == nil {
			t.Fatalf("expected an error: %v, but got nil", MalformedReqErr)
		}

		if !strings.Contains(err.Error(), MalformedReqErr.Error()) {
			t.Errorf("expected err to contain %q, but got %v", MalformedReqErr, err)
		}
	})

	t.Run("request with extra whitespces", func(t *testing.T) {
		tests := []struct {
			key   string
			value string
		}{
			{"method", "GET"},
			{"path", "/"},
			{"version", "HTTP/1.1"},
		}
		line := []byte("GET        /       HTTP/1.1")

		ParseRequestLine(line, StubMetadata)
		for _, tt := range tests {
			value := StubMetadata[tt.key]
			if value != tt.value {
				t.Errorf("got %q, want %q", value, tt.value)
			}
		}
	})
}
