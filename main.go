package main

import (
	"fmt"
	"io"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		buff := make([]byte, 8)
		var line string = ""
		for {
			n, err := f.Read(buff)
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
	}()

	return out
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

}
