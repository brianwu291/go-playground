package queue

import (
	"fmt"
	"sync"
)

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type Queue[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
	mu   sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(value T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	newNode := &Node[T]{Value: value}
	if q.head == nil {
		q.head = newNode
		q.tail = newNode
		q.size += 1
		return
	}
	q.tail.Next = newNode
	q.tail = newNode
	q.size += 1
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var zero T
	if q.head == nil {
		return zero, false
	}
	val := q.head.Value
	q.head = q.head.Next
	if q.head == nil {
		q.tail = nil
	}
	q.size -= 1
	return val, true
}

func (q *Queue[T]) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.size
}

func Demo() {
	queue := NewQueue[int]()

	queue.Enqueue(1)
	queue.Enqueue(1)

	fmt.Printf("size: %d\n", queue.Size())

	val, ok := queue.Dequeue()
	if ok {
		fmt.Printf("dequeue: %d\n", val)
	}

	fmt.Printf("size: %d\n", queue.Size())

	val, ok = queue.Dequeue()
	if ok {
		fmt.Printf("dequeue: %d\n", val)
	}

	fmt.Printf("size: %d\n", queue.Size())

	val, ok = queue.Dequeue()
	if ok {
		fmt.Printf("dequeue: %d\n", val)
	} else {
		fmt.Printf("dequeue nothing\n")
	}

	fmt.Printf("size: %d\n", queue.Size())
}
