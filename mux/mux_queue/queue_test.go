package mux_queue

import (
	"fmt"
	"log"
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


func TestNewQueue(t *testing.T) {
	var queue = NewQueue()
	queue.Push("12345")
	queue.Push(12344)
	for {
		v,err := queue.Pop()
		if err != nil{
			break
		}
		log.Print(v)
		log.Printf(" queue size = %d \n", queue.len)
	}
}
