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
		Value int
	}
}

var MalformedReqErr = errors.New("malformed request")
var ContentLengthErr = errors.New("cannot parse content-length header")

func (r *ReqMetadata) ParseRequestLine(line []byte) error {
	reqLineParts := bytes.Fields(line)
	if len(reqLineParts) < 3 {
		return fmt.Errorf("%w - response status line", MalformedReqErr)
	}

	r.Method = string(reqLineParts[0])
	r.Path = string(reqLineParts[1])
	r.Version = string(reqLineParts[2])
	return nil
}

func (r *ReqMetadata) ParseHeaders(line []byte) error {
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
			contentLength, err := strconv.Atoi(string(value))
			if err != nil {
				return fmt.Errorf("%w: %w: %v", MalformedReqErr, ContentLengthErr, err)
			}
			if contentLength < 0 {
				return fmt.Errorf("%w: %w: %v", MalformedReqErr, ContentLengthErr, err)
			}
			r.ContentLength.Value = contentLength
			r.ContentLength.isSet = true
		}
	}
	return nil
}
