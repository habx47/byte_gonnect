package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/habx47/byte_gonnect/endpoints"
)

func TestRequestLine(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		var stubMetadata ReqMetadata
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/": endpoints.GETRsp_root,
		}
		line := []byte("GET / HTTP/1.1")
		tests := []struct {
			fieldName string
			assertion func() bool
			got       func() string
			want      string
		}{
			{
				fieldName: "Method",
				assertion: func() bool { return stubMetadata.Method == "GET" },
				got:       func() string { return stubMetadata.Method },
				want:      "GET",
			},
			{
				fieldName: "Path",
				assertion: func() bool { return stubMetadata.Path == "/" },
				got:       func() string { return stubMetadata.Path },
				want:      "/",
			},
			{
				fieldName: "Version",
				assertion: func() bool { return stubMetadata.Version == "HTTP/1.1" },
				got:       func() string { return stubMetadata.Version },
				want:      "HTTP/1.1",
			},
		}

		err := stubMetadata.ParseRequestLine(line, stubRoutes, stubEndpointResps)

		if err != nil {
			t.Fatalf("failed to parser request status line: %v", err)
		}

		for _, tt := range tests {
			if !tt.assertion() {
				t.Errorf("got %v, expected %v as value for field: %v", tt.got(), tt.want, tt.fieldName)
			}
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		var stubMetadata ReqMetadata
		line := []byte("GET /")
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/": endpoints.GETRsp_root,
		}

		err := stubMetadata.ParseRequestLine(line, stubRoutes, stubEndpointResps)

		if err == nil {
			t.Fatalf("expected an error: %v, but got nil", MalformedReqErr)
		}

		if !strings.Contains(err.Error(), MalformedReqErr.Error()) {
			t.Errorf("expected err to contain %q, but got %v", MalformedReqErr, err)
		}
	})

	t.Run("request with extra whitespces", func(t *testing.T) {
		var stubMetadata ReqMetadata
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/": endpoints.GETRsp_root,
		}
		line := []byte("GET        /       HTTP/1.1")
		tests := []struct {
			fieldName string
			assertion func() bool
			got       string
			want      string
		}{
			{
				fieldName: "Method",
				assertion: func() bool { return stubMetadata.Method == "GET" },
				got:       stubMetadata.Method,
				want:      "GET",
			},
			{
				fieldName: "Path",
				assertion: func() bool { return stubMetadata.Path == "/" },
				got:       stubMetadata.Path,
				want:      "/",
			},
			{
				fieldName: "Version",
				assertion: func() bool { return stubMetadata.Version == "HTTP/1.1" },
				got:       stubMetadata.Version,
				want:      "HTTP/1.1",
			},
		}

		err := stubMetadata.ParseRequestLine(line, stubRoutes, stubEndpointResps)
		if err != nil {
			t.Fatalf("expected no err, but got: %v", err)
		}

		for _, tt := range tests {
			if !tt.assertion() {
				t.Errorf("got %v, expected %v for field %v", tt.got, tt.want, tt.fieldName)
			}
		}
	})

	t.Run("invalid version", func(t *testing.T) {
		var stubMetadata ReqMetadata
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/": endpoints.GETRsp_root,
		}
		line := []byte("GET / HTTP/2.0")

		err := stubMetadata.ParseRequestLine(line, stubRoutes, stubEndpointResps)
		if err == nil {
			t.Fatalf("expected an err but got nil")
		}

		if !errors.Is(err, HTTPVersionErr) {
			t.Errorf("got: %v, want: %v", err, HTTPVersionErr)
		}
	})

	t.Run("unsupported method", func(t *testing.T) {
		var stubMetadata ReqMetadata
		stubRoutes := map[string]bool{"/": true}
		stubEndpointResps := map[string][]byte{
			"GET/": endpoints.GETRsp_root,
		}
		line := []byte("HEAD / HTTP/1.1")

		err := stubMetadata.ParseRequestLine(line, stubRoutes, stubEndpointResps)
		if err == nil {
			t.Fatalf("expected an err but got nil")
		}

		if !errors.Is(err, MethodNotSupportedErr) {
			t.Errorf("got: %v, want: %v", err, MethodNotSupportedErr)
		}
	})
}

func TestParseHeaders(t *testing.T) {
	t.Run("valid field", func(t *testing.T) {
		var stubMetadata ReqMetadata

		line := []byte("Content-Length: 10")
		err := stubMetadata.ParseHeaders(line)
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}

		got := stubMetadata.ContentLength.Value
		want := uint64(10)
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		var stubMetadata ReqMetadata
		line := []byte("Content-Length: abc")

		err := stubMetadata.ParseHeaders(line)
		if err == nil {
			t.Fatalf("expected error: %v, but got nil", ContentLengthErr)
		}

		if !errors.Is(err, ContentLengthErr) {
			t.Errorf("got %v, want %v", err, ContentLengthErr)
		}
	})
}
