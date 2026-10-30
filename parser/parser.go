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

var MalformedReqErr = errors.New("malformed request")
var ContentLengthErr = errors.New("cannot parse content-length header")
var MalformedHeaderErr = errors.New("malformed header")
var HTTPVersionErr = errors.New("unsupported http version")
var MethodNotSupportedErr = errors.New("method not currently supported")
var NotFoundErr = errors.New("not found")
var MethodNotAllowedErr = errors.New("method not allowed")

func (r *ReqMetadata) ParseRequestLine(line []byte, routes endpoints.Routes, endpointResps endpoints.EndpointResps) error {
	reqLineParts := bytes.Fields(line)
	if len(reqLineParts) != 3 {
		return fmt.Errorf("%w - response status line", MalformedReqErr)
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

	if _, ok := routes[path]; !ok {
		return fmt.Errorf("%w", NotFoundErr)
	}

	endpoint := method + path
	if _, ok := endpointResps[endpoint]; !ok {
		return fmt.Errorf("%w", MethodNotAllowedErr)
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

	if found {
		key = bytes.TrimSpace(key)
		value = bytes.TrimSpace(value)
	}

	if bytes.EqualFold(key, []byte("Content-Length")) {
		if r.ContentLength.isSet {
			return fmt.Errorf("%v: multiple Content-Length headers in the request", MalformedReqErr)
		}

		contentLength, err := strconv.ParseUint(string(value), 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %w: %v", MalformedReqErr, ContentLengthErr, err)
		}
		r.ContentLength.Value = contentLength
		r.ContentLength.isSet = true

	}
	return nil
}
