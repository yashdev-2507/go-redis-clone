package main

import (
	"hash/fnv"
	"sync"
	"time"
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
	Type      string
	Str       string
	List      *LinkedList
	Set       map[string]struct{}
	Hash      map[string]string
	ExpiresAt time.Time
}

type Shard struct {
	mu   sync.RWMutex
	data map[string]*Value
}

type ShardedStore struct {
	shards [16]*Shard
}

func startupShardedStore() *ShardedStore {
	s := &ShardedStore{}
	for i := range s.shards {
		s.shards[i] = &Shard{data: make(map[string]*Value)}
	}
	return s

}

func getShardIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() % 16)
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

func (s *ShardedStore) Get(key string) (*Value, bool) {
	ind := getShardIndex(key)
	s.shards[ind].mu.RLock()
	defer s.shards[ind].mu.RUnlock()
	value, ok := s.shards[ind].data[key]
	return value, ok
}

func (s *ShardedStore) Delete(key string) {
	ind := getShardIndex(key)
	s.shards[ind].mu.Lock()
	defer s.shards[ind].mu.Unlock()

	delete(s.shards[ind].data, key)
}

func (s *ShardedStore) Set(key string, value *Value) {
	ind := getShardIndex(key)
	s.shards[ind].mu.Lock()
	defer s.shards[ind].mu.Unlock()
	s.shards[ind].data[key] = value
}
