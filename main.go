package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
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

func parseRESP(reader *bufio.Reader) ([]string, error) {
	s, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	str := s[1 : len(s)-2]         // strip the leading '*' and trailing '\r\n', leaving just the digit characters (e.g. "*3\r\n" -> "3")
	size, err := strconv.Atoi(str) // convert the digit string into an actual int we can use
	if err != nil {
		return nil, err
	}
	var args []string
	for i := 0; i < size; i++ {
		s, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		str := s[1 : len(s)-2]

		byte_size, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, byte_size)
		_, err_1 := io.ReadFull(reader, buf)
		if err_1 != nil {
			return nil, err_1
		}

		args = append(args, string(buf))

		_, err_2 := reader.ReadString('\n')
		if err_2 != nil {
			return nil, err_2
		}

	}
	return args, nil
}
