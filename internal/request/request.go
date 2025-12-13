package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Status int

const (
	INITIALIZED Status = iota
	DONE
)

type Request struct {
	RequestLine  RequestLine
	ParserStatus Status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func checkUpperCase(str string) bool {
	for _, r := range str {
		if r < 'A' || r > 'Z' {
			return false
		}
	}

	return true
}

func parseRequestLine(req string) (*RequestLine, int, error) {
	req_line := strings.Split(req, "\r\n")[0]
	if req_line == "" {
		return nil, 0, nil
	}
	parts := strings.Split(req_line, " ")

	if len(parts) != 3 {
		return nil, 0, fmt.Errorf("invalid request line")
	}

	method, target, protocol_version := parts[0], parts[1], parts[2]

	if !checkUpperCase(method) {
		return nil, 0, fmt.Errorf("method in request line should be uppercase")
	}

	version := strings.Split(protocol_version, "/")[1]
	if version != "1.1" {
		return nil, 0, fmt.Errorf("%s: incorrect protocol version. 1.1 supported", version)
	}

	res := RequestLine{HttpVersion: version, RequestTarget: target, Method: method}
	return &res, len(req_line), nil
}

func (r *Request) parse(data []byte) (int, error) {
	r.ParserStatus = INITIALIZED

	requestLine, parsedBytes, err2 := parseRequestLine(string(data))
	if err2 != nil {
		return 0, errors.Join(fmt.Errorf("some error occured while parsing line: "), err2)
	}

	r.RequestLine = *requestLine
	return parsedBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := Request{}
	buff := make([]byte, 1024)
	for {
		req, err := io.ReadAtLeast(reader, buff, 8)
		if err != nil {
			return nil, err
		}

	}
	return &request, nil
}
