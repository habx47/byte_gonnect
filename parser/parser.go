package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
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
var URIErr = errors.New("URI too long")
var ContentLengthErr = errors.New("cannot parse content-length header")
var HeaderErr = errors.New("headers too large")

func (r *ReqMetadata) ParseRequestLine(line []byte, isPrefix bool) error {
	if isPrefix {
		return fmt.Errorf("%w - %w - response status line", MalformedReqErr, URIErr)
	}

	reqLineParts := bytes.Fields(line)
	if len(reqLineParts) != 3 {
		return fmt.Errorf("%w - response status line", MalformedReqErr)
	}

	r.Method = string(reqLineParts[0])
	r.Path = string(reqLineParts[1])
	r.Version = string(reqLineParts[2])
	return nil
}

func (r *ReqMetadata) ParseHeaders(line []byte, isPrefix bool) error {
	if isPrefix {
		return fmt.Errorf("%w: %w", MalformedReqErr, HeaderErr)
	}

	key, value, found := bytes.Cut(line, []byte(":"))

	if !found {
		return fmt.Errorf("%w: %w", MalformedReqErr, ContentLengthErr)
	}

	if found {
		key = bytes.TrimSpace(key)
		value = bytes.TrimSpace(value)
	}

	if bytes.EqualFold(key, []byte("Content-Length")) {
		if !r.ContentLength.isSet {
			contentLength, err := strconv.ParseUint(string(value), 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %w: %v", MalformedReqErr, ContentLengthErr, err)
			}
			r.ContentLength.Value = contentLength
			r.ContentLength.isSet = true
		}

		if r.ContentLength.isSet {
			return fmt.Errorf("%v: multiple Content-Length headers in the request", MalformedReqErr)
		}
	}
	return nil
}
