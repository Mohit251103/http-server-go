package server

import (
	"fmt"
	"io"
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

func (s *Server) handle(conn io.WriteCloser) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!")
	conn.Write(out)
	conn.Close()
	// request, err2 := request.RequestFromReader(conn)
	// if err2 != nil {
	// 	return
	// }

	// fmt.Printf("Request line: \n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
	// fmt.Printf("Header: ")
	// for key, value := range request.Headers {
	// 	if key == "content-length" {
	// 		continue
	// 	}
	// 	fmt.Printf("- %s: %s\n", key, value)
	// }
	// fmt.Printf("Body:\n%s", request.Body)
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
