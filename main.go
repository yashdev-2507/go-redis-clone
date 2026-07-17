package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("failed to bind : ", err)
		return
	}
	defer listener.Close()
	fmt.Println("listening on : 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handleconnection(conn)
	}

}

func handleconnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := parseRESP(reader)
		if err != nil {
			fmt.Println("client disconnected or bad input:", err)
			return
		}
		fmt.Println("received command: ", command)
	}
}
