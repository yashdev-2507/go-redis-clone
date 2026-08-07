package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
)

type ListNode struct {
	Val  string
	Prev *ListNode
	Next *ListNode
}

type LinkedList struct {
	Head *ListNode
	Tail *ListNode
	Size int
}

type Value struct {
	Type string
	Str  string
	List *LinkedList
	Set  map[string]struct{}
	Hash map[string]string
}

func (ll *LinkedList) PushFront(val string) {
	node := &ListNode{Val: val}

	if ll.Head == nil {
		ll.Head = node
		ll.Tail = node
		ll.Head.Prev = nil
		ll.Tail.Next = nil

	} else {
		ll.Head.Prev = node
		node.Next = ll.Head
		ll.Head = node
		ll.Head.Prev = nil

	}
	ll.Size++
}

func (ll *LinkedList) PushBack(val string) {
	node := &ListNode{Val: val}

	if ll.Head == nil {
		ll.Head = node
		ll.Tail = node
	} else {
		ll.Tail.Next = node
		node.Prev = ll.Tail
		ll.Tail = node
	}

	ll.Size++
}

func (ll *LinkedList) PopFront() (string, bool) {
	switch ll.Size {
	case 0:
		return "", false
	case 1:
		val := ll.Head.Val
		ll.Head = nil
		ll.Tail = nil
		ll.Size--
		return val, true

	}
	val := ll.Head.Val
	ll.Head = ll.Head.Next
	ll.Head.Prev = nil
	ll.Size--
	return val, true
}

func (ll *LinkedList) PopBack() (string, bool) {
	switch ll.Size {
	case 0:
		return "", false
	case 1:
		val := ll.Tail.Val
		ll.Head = nil
		ll.Tail = nil
		ll.Size--
		return val, true
	}
	val := ll.Tail.Val
	ll.Tail = ll.Tail.Prev
	ll.Tail.Next = nil
	ll.Size--
	return val, true
}

func (ll *LinkedList) LRange(start, stop int) []string {
	result := []string{}
	temp := ll.Head
	i := 0
	if start < 0 {
		start += ll.Size
	}
	if stop < 0 {
		stop += ll.Size
	}
	for temp != nil {
		if i >= start && i <= stop {
			result = append(result, temp.Val)
		}
		temp = temp.Next
		i++
	}
	return result
}

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

func SEThandler(args []string, store map[string]*Value) string {
	if len(args) != 3 {
		return "-ERR wrong numbers of arguments for 'set' command\r\n"
	}
	store[args[1]] = &Value{Type: "string", Str: args[2]}
	return "+OK\r\n"
}
func GEThandler(args []string, store map[string]*Value) string {
	if len(args) != 2 {
		return "-ERR wrong number of arguments for 'get' command\r\n"
	}
	value, ok := store[args[1]]
	if !ok {
		return "$-1\r\n"
	}
	if value.Type != "string" {
		return "-WRONGTYPE Operation agianst a key holding the wrong kind of value\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value.Str), value.Str)

}

func LPUSHhandler(args []string, store map[string]*Value) string {
	if len(args) < 3 {
		return "-ERR wrong number of arguements	for 'LPUSH' command\r\n"
	}
	value, ok := store[args[1]]
	if !ok {
		value = &Value{Type: "list", List: &LinkedList{}}
		store[args[1]] = value

	} else if value.Type != "list" {
		return "-WRONGTYPE operation against a key holding the wrong kind of value\r\n"
	}
	i := 2
	for i != len(args) {
		value.List.PushFront(args[i])
		i++
	}

	return fmt.Sprintf(":%d\r\n", value.List.Size)
}
func RPUSHhandler(args []string,store map[string]*Value)string{
	if len(args) < 3{
		return "-ERR wrong number of arguments for rpush command\r\n"
	}
	value,ok := store[args[1]]
	if !ok {
		value = &Value{Type : "list", List : &LinkedList{}}
		store[args[1]] = value
	}else if value.Type != "list"{
		return "-WRONGTYPE operation againt a key holding the wrong  kind of value\r\n"
	}
	i:= 2
	for i != len(args){
		value.List.PushBack(args[i])
		i++;
	}
	return fmt.Sprintf(":%d\r\n", value.List.Size)

}

func RPOPhandler(args []string, store map[string]*Value) string {
	if len(args) != 2 {
		return "-ERR wrong number of arguments for 'RPOP' command\r\n"
	}
	value, ok := store[args[1]]
	if !ok {
		return "$-1\r\n"
	}
	if value.Type != "list" {
		return "-WRONGTYPE operation against a key holding the wrong kind of value\r\n"
	}
	val, ok := value.List.PopBack()
	if !ok {
		return "$-1\r\n"
	}
	if value.List.Size == 0 {
		delete(store, args[1])
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
}

func LPOPhandler(args []string, store map[string]*Value) string {
	if len(args) != 2 {
		return "-ERR wrong number of arguments for 'LPOP' command\r\n"
	}
	value, ok := store[args[1]]
	if !ok {
		return "$-1\r\n"
	}
	if value.Type != "list" {
		return "-WRONGTYPE operation against a key holding the wrong kind of value\r\n"
	}
	val, ok := value.List.PopFront()
	if !ok {
		return "$-1\r\n"
	}
	if value.List.Size == 0 {
		delete(store, args[1])
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
}

func LRANGEhandler(args []string, store map[string]*Value) string {
	if len(args) != 4 {
		return "-ERR wrong number of arguments for 'LRANGE' command\r\n"
	}
	value, ok := store[args[1]]
	if !ok {
		return "*0\r\n"
	}
	if value.Type != "list" {
		return "-WRONGTYPE operation against a key holding the wrong kind of value\r\n"
	}
	start, err := strconv.Atoi(args[2])
	if err != nil {
		return "-ERR value is not an integer or out of range\r\n"
	}
	stop, err := strconv.Atoi(args[3])
	if err != nil {
		return "-ERR value is not an integer or out of range\r\n"
	}

	result := value.List.LRange(start, stop)

	
	
	result_string := fmt.Sprintf("*%d\r\n", len(result))
	for _,val := range result {
		result_string += fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
	}
	return result_string
}

func dispatch(args []string, store map[string]*Value) string {
	if len(args) == 0 {
		return "-ERR you have entered nothing\r\n"
	}
	switch args[0] {
	case "set":
		return SEThandler(args, store)
	case "get":
		return GEThandler(args, store)
	case "lpush":
		return LPUSHhandler(args, store)
	case "rpush":
		return RPUSHhandler(args, store)
	case "lpop" :
		return LPOPhandler(args,store)
	case "rpop":
		return RPOPhandler(args,store)
	case "lrange":
		return LRANGEhandler(args,store)
	default:
		return "-ERR wrong input\r\n"
	}

}
