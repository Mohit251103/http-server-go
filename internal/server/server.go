package server

import (
	"bytes"
	"fmt"
	"http-server/internal/request"
	"http-server/internal/response"
	"io"
	"net"
	"strconv"
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
	// parsing request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Error", err)
	}

	// parsing response
	res := bytes.NewBuffer(make([]byte, 1024))

	herr := s.handler(res, req)
	if herr != nil {
		_ = WriteHandlerError(res, herr)
	} else {
		err := response.WriteStatusLine(res, 200)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		contentLength, err := strconv.Atoi(req.Headers["content-length"])
		if err != nil {
			contentLength = 0
		}
		headers := response.GetDefaultHeaders(contentLength)
		err = response.WriteHeaders(res, headers)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	fmt.Println(string(res.Bytes()))
	// if err != nil {
	// 	fmt.Println("Error writing response to connection:", err)
	// }

	err = conn.Close()
	if err != nil {
		fmt.Println("Error closing connection:", err)
	}
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
