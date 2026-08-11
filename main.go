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
	store := make(map[string]*Value)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handleconnection(conn, store)
	}

}

func handleconnection(conn net.Conn, store map[string]*Value) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := parseRESP(reader)
		if err != nil {
			fmt.Println("client disconnected or bad input:", err)
			return
		}

		reply := dispatch(command, store)
		conn.Write([]byte(reply))
	}
}
