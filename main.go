package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func getLinesChannel(conn net.Conn) <-chan string {

	out := make(chan string, 1)

	go func() {
		defer close(out)

		buff := make([]byte, 8)
		var line string = ""
		for {
			n, err := conn.Read(buff)
			buff = buff[:n]
			if err == nil {
				start := 0
				for i, b := range buff {
					if b == '\n' {
						line += string(buff[start:i])
						out <- line
						line = ""
						start = i + 1
					}
				}

				line += string(buff[start:n])
			}

			if err == io.EOF {
				if len(line) != 0 {
					out <- line
				}
				break
			}
		}
		err := conn.Close()
		if err == nil {
			fmt.Println("Connection has been closed!")
		}
	}()

	return out
}

func main() {
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("error listening: %s\n", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err == nil {
			fmt.Println("Connection has been accepted!")
		}
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
	}

}
