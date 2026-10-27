package parser

import (
	"bytes"
	"errors"
	"fmt"
)

type ReqMetadata map[string]string

var MalformedReqErr = errors.New("malformed request")

func ParseRequestLine(line []byte, metadata ReqMetadata) error {
	reqLineParts := bytes.Fields(line)
	if len(reqLineParts) < 3 {
		return fmt.Errorf("%w - response status line", MalformedReqErr)
	}

	metadata["method"] = string(reqLineParts[0])
	metadata["path"] = string(reqLineParts[1])
	metadata["version"] = string(reqLineParts[2])
	return nil
}

func ParseHeaders(line []byte) {
}
