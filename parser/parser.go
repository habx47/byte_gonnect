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
	ContentLength int
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
	if key, value, found := bytes.Cut(line, []byte(":")); found {
		key = bytes.TrimSpace(key)
		value = bytes.TrimSpace(value)

		if bytes.EqualFold(key, []byte("Content-Length")) {
			contentLength, err := strconv.Atoi(string(value))
			if err != nil {
				return fmt.Errorf("%w: %v", ContentLengthErr, err)
			}
			r.ContentLength = contentLength
		}
	}
	return nil
}
