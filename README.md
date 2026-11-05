# # byte_gonnect

A minimal TCP server written in Go that implements a small subset of HTTP/1.1 request handling.

The project explores network communication, request parsing, connection management, error handling, and minimal resource limits using Go's standard library.

## Features

- TCP listener and concurrent connection handling using goroutines.
- Manual HTTP request-line and header parsing.
- Request validation, including methods, HTTP version, header fields, and `Content-Length`.
- Basic routing based on HTTP method and request path.
- Request-body reading based on the declared `Content-Length`.
- HTTP error responses for malformed and oversized requests.
- Connection cleanup and read deadlines.
- Tests for successful requests and selected read, write, and parsing failures.

## Design

The server uses Go's `net` package for TCP communication and `bufio` for buffered request reading.

Each accepted connection is handed to a separate goroutine. The connection handler owns the connection lifecycle: it reads and validates one request, routes it, writes the response, and closes the connection.

Request parsing is separated from routing and endpoint response handling. This keeps protocol-level concerns distinct from application behavior.

The server handles one request per connection and closes it after responding. Persistent connections and chunked transfer encoding are not implemented.
