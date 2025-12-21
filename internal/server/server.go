package server

import (
	"bytes"
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"net"
)

type Server struct {
	listener net.Listener
	handler  Handler
}

func Serve(port int, handlerFunc Handler) (*Server, error) {
	netAddress := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", netAddress)

	if err != nil {
		return nil, err
	}

	return &Server{listener: listener, handler: handlerFunc}, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()
	// parsing request
	headers := response.GetDefaultHeaders(0)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		WriteHandlerError(conn, hErr)
		return
	}

	// parsing response
	res := bytes.NewBuffer([]byte{})
	var body []byte = nil
	statusCode := response.StatusOk
	herr := s.handler(res, req)
	if herr != nil {
		statusCode = herr.StatusCode
		body = []byte(herr.Message)
	} else {
		body = res.Bytes()
	}

	headers["Content-Length"] = fmt.Sprintf("%d", len(body))
	response.WriteStatusLine(conn, statusCode)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
	// out := []byte("HTTP/1.1 200 OK\r\nContent-Length: 12\r\nConnection: close\r\nContent-Type: text/plain\r\n\r\nHello world!")
	// conn.Write(out)
	// conn.Close()
}

func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}

		go s.handle(conn)
	}
}
