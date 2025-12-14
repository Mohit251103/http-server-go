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
	idx := strings.Index(req, "\r\n")

	if idx == -1 {
		return nil, 0, nil
	}

	req_line := strings.Split(req, "\r\n")[0]
	parts := strings.Split(req_line, " ")

	method, target, version := parts[0], parts[1], parts[2]

	if !checkUpperCase(method) {
		return nil, 0, fmt.Errorf("method in request line should be uppercase")
	}

	// version := strings.Split(protocol_version, "/")[1]
	if version != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("%s: incorrect protocol version. 1.1 supported", version)
	}

	res := RequestLine{HttpVersion: "1.1", RequestTarget: target, Method: method}
	return &res, len(req_line), nil
}

func (r *Request) parse(data []byte) (int, error) {

	requestLine, parsedBytes, err2 := parseRequestLine(string(data))
	if err2 != nil {
		return 0, errors.Join(fmt.Errorf("some error occured while parsing line: "), err2)
	}

	if parsedBytes == 0 {
		return 0, nil
	}

	r.ParserStatus = DONE
	if requestLine != nil {
		r.RequestLine = *requestLine
	}
	return parsedBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := Request{}
	buff := make([]byte, 1024)
	bufLen := 0
	for request.ParserStatus != DONE {
		n, err := reader.Read(buff[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readBytes, err2 := request.parse(buff[:bufLen])
		if err2 != nil {
			return nil, err2
		}

		copy(buff, buff[readBytes:bufLen])
		bufLen -= readBytes
	}
	return &request, nil
}
