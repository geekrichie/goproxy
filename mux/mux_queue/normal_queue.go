package mux_queue

import (
	"errors"
	"unsafe"
)

type Queue struct{
	Len    int
	head  unsafe.Pointer
	tail unsafe.Pointer
}

type Item struct {
	value interface{}
	next unsafe.Pointer
}

func NewQueue() *Queue{
	var item = Item{}
	return &Queue{
		head: unsafe.Pointer(&item),
		tail: unsafe.Pointer(&item),
		Len : 0,
	}
}

func (q *Queue) Push (value interface{}) {
	var item = Item{value:value, next: nil}
	tail := (*Item)(q.tail)
	tail.next = unsafe.Pointer(&item)
	q.tail = tail.next
	q.Len = q.Len + 1
}

func (q *Queue)  Pop() (interface{}, error){
	if q.Len == 0 {
		return nil,errors.New("the queue is empty")
	}
	head  := (*Item)(q.head)
	next := head.next
	var item = (*Item)(next)
	q.head = next
	q.Len = q.Len - 1
	return item.value,nil
}