package server

import (
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteHandlerError(w io.Writer, err *HandlerError) error {
	_, e := w.Write([]byte(fmt.Sprintf("error: %d\r\nmessage: %s", err.StatusCode, err.Message)))

	if e != nil {
		return e
	}

	return nil
}
