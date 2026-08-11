package main

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
