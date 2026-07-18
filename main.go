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
	store := make(map[string]string)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handleconnection(conn,store)
	}

}

func handleconnection(conn net.Conn, store map[string]string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		command, err := parseRESP(reader)
		if err != nil {
			fmt.Println("client disconnected or bad input:", err)
			return
		}

		reply := dispatch(command, store);
		conn.Write([]byte(reply));
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

func SEThandler(args []string, store map[string]string) string{
	if len(args) != 3 {
		return "-ERR wrong numbers of arguments for 'set' command\r\n"
	}
	store[args[1]] = args[2]
	return "+OK\r\n"
}
func GEThandler(args []string, store map[string]string)string {
	if len(args) != 2{
		return "-ERR wrong number of arguments for 'get' command\r\n"
	}
	value,ok := store[args[1]]
	if ok {
		return fmt.Sprintf("$%d\r\n%s\r\n", len(value),value)
	}else{
		return "$-1\r\n"
	}

}

func dispatch(args []string, store map[string]string) string{
	if len(args) == 0{
		return "-ERR you have entered nothing\r\n"
	}
	switch args[0] {
	case "set":
		return SEThandler(args, store)
	case "get":
		return GEThandler(args, store)
	default:
		return "-ERR wrong input\r\n"
	}
}
