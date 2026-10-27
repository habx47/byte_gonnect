package parser

import (
	"strings"
	"testing"
)

func TestRequestLine(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		var stubMetadata ReqMetadata
		line := []byte("GET / HTTP/1.1")
		tests := []struct {
			feildName string
			assertion func() bool
			got       string
			want      string
		}{
			{
				feildName: "Method",
				assertion: func() bool { return stubMetadata.Method == "GET" },
				got:       stubMetadata.Method,
				want:      "GET",
			},
			{
				feildName: "Path",
				assertion: func() bool { return stubMetadata.Path == "/" },
				got:       stubMetadata.Path,
				want:      "/",
			},
			{
				feildName: "Version",
				assertion: func() bool { return stubMetadata.Version == "HTTP/1.1" },
				got:       stubMetadata.Version,
				want:      "HTTP/1.1",
			},
		}

		err := stubMetadata.ParseRequestLine(line)

		if err != nil {
			t.Fatalf("failed to parser request status line: %v", err)
		}

		for _, tt := range tests {
			if !tt.assertion() {
				t.Errorf("got %v, expected %v", tt.got, tt.want)
			}
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		var stubMetadata ReqMetadata
		line := []byte("GET /")

		err := stubMetadata.ParseRequestLine(line)

		if err == nil {
			t.Fatalf("expected an error: %v, but got nil", MalformedReqErr)
		}

		if !strings.Contains(err.Error(), MalformedReqErr.Error()) {
			t.Errorf("expected err to contain %q, but got %v", MalformedReqErr, err)
		}
	})

	t.Run("request with extra whitespces", func(t *testing.T) {
		var stubMetadata ReqMetadata
		line := []byte("GET        /       HTTP/1.1")
		tests := []struct {
			feildName string
			assertion func() bool
			got       string
			want      string
		}{
			{
				feildName: "Method",
				assertion: func() bool { return stubMetadata.Method == "GET" },
				got:       stubMetadata.Method,
				want:      "GET",
			},
			{
				feildName: "Path",
				assertion: func() bool { return stubMetadata.Path == "/" },
				got:       stubMetadata.Path,
				want:      "/",
			},
			{
				feildName: "Version",
				assertion: func() bool { return stubMetadata.Version == "HTTP/1.1" },
				got:       stubMetadata.Version,
				want:      "HTTP/1.1",
			},
		}

		stubMetadata.ParseRequestLine(line)

		for _, tt := range tests {
			if !tt.assertion() {
				t.Errorf("got %v, expected %v", tt.got, tt.want)
			}
		}
	})
}

func TestParseHeaders(t *testing.T) {
	var stubMetadata ReqMetadata
	line := []byte("Content-Length: 10")
	stubMetadata.ParseHeaders(line)

	got := stubMetadata.ContentLength
	want := 10
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
