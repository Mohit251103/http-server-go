package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	r, _ := net.ResolveUDPAddr("udp", ":42069")
	conn, _ := net.DialUDP("udp", nil, r)
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		lines, err := reader.ReadString('\n')
		if err != nil {
			panic("Problem while trying to read line from bufio.Reader")
		}

		_, err2 := conn.Write([]byte(lines))
		if err2 != nil {
			panic("Error while writing to udp")
		}
	}
}
