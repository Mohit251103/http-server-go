package server

import (
	"fmt"
	"http-server/internal/request"
	"net"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	netAddress := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", netAddress)

	if err != nil {
		return nil, err
	}

	return &Server{listener: listener}, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) handle(conn net.Conn) {
	request, err2 := request.RequestFromReader(conn)
	if err2 != nil {
		return
	}

	fmt.Printf("Request line: \n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
	fmt.Printf("Header: ")
	for key, value := range request.Headers {
		fmt.Printf("- %s: %s\n", key, value)
	}
	fmt.Printf("Body:\n%s", request.Body)
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
