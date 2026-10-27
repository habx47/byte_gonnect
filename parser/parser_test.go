package parser

import "testing"

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
				t.Errorf("expected key: %v to exist", tt.key)
			}
		}
	})
}
