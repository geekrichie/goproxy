package mux_queue

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkEnQueue(b *testing.B) {
	queue := NewLKQueue()
	var waiter sync.WaitGroup
	for i := 0; i< 100;i++{
		waiter.Add(1)
		go func(i int){
			queue.Enqueue(i)
			waiter.Done()
		}(i)
	}
	waiter.Wait()
	for i := 0 ; i< 100 ;i ++ {
		fmt.Println(queue.Dequeue())
	}
}

func TestLKQueue_Enqueue(t *testing.T) {
	queue := NewLKQueue()
	var waiter sync.WaitGroup
	for i := 0; i< 100;i++{
		waiter.Add(1)
		go func(i int){
			queue.Enqueue(i)
			waiter.Done()
		}(i)
	}
	waiter.Wait()
	for i := 0 ; i< 100 ;i ++ {
		fmt.Println(queue.Dequeue())
	}
}


type ListElement struct{
	buf    []byte
	L      int
	isPart bool
}
func TestQueue(t *testing.T) {
	bufQueue := NewLKQueue()
	b := []byte("123")
	listelem := new(ListElement)
	listelem.buf = b
	listelem.L = len(b)
	bufQueue.Enqueue(listelem)
	r := bufQueue.Dequeue()
	fmt.Println(r)
}
