package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
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

func parseRequestLine(req string) (*RequestLine, error) {
	req_line := strings.Split(req, "\r\n")[0]
	parts := strings.Split(req_line, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	method, target, protocol_version := parts[0], parts[1], parts[2]

	if !checkUpperCase(method) {
		return nil, fmt.Errorf("method in request line should be uppercase")
	}

	version := strings.Split(protocol_version, "/")[1]
	if version != "1.1" {
		return nil, fmt.Errorf("%s: incorrect protocol version. 1.1 supported", version)
	}

	res := RequestLine{HttpVersion: version, RequestTarget: target, Method: method}
	return &res, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	req, err := io.ReadAll(reader)
	if err != nil {
		panic(err.Error())
	}

	requestLine, err2 := parseRequestLine(string(req))
	if err2 != nil {
		return nil, err2
	}
	request := Request{RequestLine: *requestLine}
	return &request, nil
}
