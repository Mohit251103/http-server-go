package main

import (
	"http-server/internal/request"
	"http-server/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	handlerFunc := func(w io.Writer, req *request.Request) *server.HandlerError {
		path := req.RequestLine.RequestTarget
		switch path {
		case "/yourproblem":
			return &server.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
		case "/myproblem":
			return &server.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
		default:
			w.Write([]byte("All good, frfr\n"))
		}
		return nil
	}
	server, err := server.Serve(port, handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server started on port", port)
	server.Listen()
	defer server.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
