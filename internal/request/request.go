package request

import (
	"errors"
	"fmt"
	"http-server/internal/headers"
	"io"
	"strings"
)

type Status string

const (
	INITIALIZED                Status = "init"
	DONE                       Status = "Done"
	requestStateParsingHeaders        = "Parsing Headers"
)

type Request struct {
	RequestLine  RequestLine
	ParserStatus Status
	Headers      headers.Headers
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
	return &res, len(req_line + "\r\n"), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.ParserStatus == requestStateParsingHeaders {
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.ParserStatus = DONE
		}
		return n, nil
	}

	requestLine, parsedBytes, err2 := parseRequestLine(string(data))
	if err2 != nil {
		return 0, errors.Join(fmt.Errorf("some error occured while parsing line: "), err2)
	}

	if parsedBytes == 0 {
		return 0, nil
	}

	if requestLine != nil {
		r.RequestLine = *requestLine
	}

	if r.ParserStatus == INITIALIZED {
		r.ParserStatus = requestStateParsingHeaders
		return parsedBytes, nil
	}

	return parsedBytes, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := Request{ParserStatus: INITIALIZED, Headers: headers.NewHeaders()}
	buff := make([]byte, 1024)
	bufLen := 0
	for request.ParserStatus != DONE {
		n, err := reader.Read(buff[bufLen:])
		if err != nil {
			if err == io.EOF {
				request.ParserStatus = DONE
				break
			}
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
