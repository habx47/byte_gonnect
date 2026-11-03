package parser

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/habx47/byte_gonnect/endpoints"
)

type ReqMetadata struct {
	Method        string
	Path          string
	Version       string
	ContentLength struct {
		isSet bool
		Value uint64
	}
}

var BadRequestErr = errors.New("bad request")
var ContentLengthErr = errors.New("cannot parse content-length header")
var MalformedHeaderErr = errors.New("malformed header")
var EmptyHeaderK_VErr = errors.New("header key or value is empty")
var WhitespaceErr = errors.New("extra whitespace in header key")
var InvalidHeaderKeyToken = errors.New("invalid token in header key")
var InvalidHeaderValueToken = errors.New("invalid token in header value")
var HTTPVersionErr = errors.New("unsupported http version")
var MethodNotSupportedErr = errors.New("method not currently supported")

func (r *ReqMetadata) ParseRequestLine(line []byte) error {
	reqLineParts := bytes.Fields(line)
	if len(reqLineParts) != 3 {
		return fmt.Errorf("%w - response status line", BadRequestErr)
	}

	method := string(reqLineParts[0])
	path := string(reqLineParts[1])
	version := string(reqLineParts[2])

	if !strings.EqualFold(version, endpoints.HTTPVersion) {
		return fmt.Errorf("%w - response status line", HTTPVersionErr)
	}

	if !slices.Contains(endpoints.SupportedMethods, method) {
		return fmt.Errorf("%w - response status line", MethodNotSupportedErr)
	}

	r.Method = method
	r.Path = path
	r.Version = version
	return nil
}

func (r *ReqMetadata) ParseHeaders(line []byte) error {
	key, value, found := bytes.Cut(line, []byte(":"))

	if !found {
		return fmt.Errorf("%w", MalformedHeaderErr)
	}
	value = bytes.TrimSpace(value)

	if len(bytes.TrimSpace(key)) == 0 {
		return fmt.Errorf("%w %w", MalformedHeaderErr, EmptyHeaderK_VErr)
	}

	if bytes.ContainsAny(key, " ") {
		return fmt.Errorf("%w %w", MalformedHeaderErr, WhitespaceErr)
	}

	for i := range key {
		if !endpoints.ValidTokenChars[key[i]] {
			return fmt.Errorf("%w %w", MalformedHeaderErr, InvalidHeaderKeyToken)
		}
	}

	for _, char := range value {
		if (char < 0x20 || char > 0x7E) && char != 0x09 {
			return fmt.Errorf("%w %w", MalformedHeaderErr, InvalidHeaderValueToken)
		}
	}

	if bytes.EqualFold(key, []byte("Content-Length")) && len(value) != 0 {
		if r.ContentLength.isSet {
			return fmt.Errorf("%v: multiple Content-Length headers in the request", BadRequestErr)
		}

		if len(value) > 20 {
			return fmt.Errorf("%w: content-length value overflow", ContentLengthErr)
		}

		for i := range value {
			if value[i] < '0' || value[i] > '9' {
				return fmt.Errorf("%w %w", ContentLengthErr, InvalidHeaderValueToken)
			}
		}

		contentLength, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w: %v", BadRequestErr, ContentLengthErr, err)
		}
		r.ContentLength.Value = contentLength
		r.ContentLength.isSet = true

	} else {
		return fmt.Errorf("%w %w", ContentLengthErr, EmptyHeaderK_VErr)
	}
	return nil
}
