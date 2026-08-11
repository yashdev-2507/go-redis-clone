package main

import (
	"fmt"
	"strconv"
)

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
		return "-ERR wrong number of arguments for 'LPUSH' command\r\n"
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

func RPUSHhandler(args []string, store map[string]*Value) string {
	if len(args) < 3 {
		return "-ERR wrong number of arguments for rpush command\r\n"
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
		value.List.PushBack(args[i])
		i++
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
	for _, val := range result {
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
	case "lpop":
		return LPOPhandler(args, store)
	case "rpop":
		return RPOPhandler(args, store)
	case "lrange":
		return LRANGEhandler(args, store)
	default:
		return "-ERR wrong input\r\n"
	}

}
