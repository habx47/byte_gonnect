package parser

import (
	"bytes"
	"errors"
	"fmt"
)

type ReqMetadata map[string]string

var MalformedReqErr = errors.New("malformed request")

func ParseRequestLine(line []byte, metadata ReqMetadata) error {
	reqLineParts := bytes.Split(line, []byte(" "))
	if len(reqLineParts) < 3 {
		return fmt.Errorf("%w - response status line", MalformedReqErr)
	}

	metadata["method"] = string(bytes.TrimSpace(reqLineParts[0]))
	metadata["path"] = string(bytes.TrimSpace(reqLineParts[1]))
	metadata["version"] = string(bytes.TrimSpace(reqLineParts[2]))
	return nil
}

func ParseHeaders(line []byte) {
}
