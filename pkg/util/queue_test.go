package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueue(t *testing.T) {
	assertion := assert.New(t)
	queue := NewQueue[int]()
	assertion.Equal(0, queue.Len())
	// Insert test
	for i := 0; i < 100; i++ {
		queue.Enqueue(i)
	}
	assertion.Equal(100, queue.Len())
	{
		expectedQueue := make([]int, 100)
		for i := 0; i < 100; i++ {
			expectedQueue[i] = i
		}
		assertion.Equal(expectedQueue, queue.CopyAsArray())
	}
	// Remove some elements
	for i := 0; i < 50; i++ {
		assertion.Equal(i, queue.Dequeue())
	}
	assertion.Equal(50, queue.Len())
	// Again add elements
	for i := 0; i < 100; i++ {
		queue.Enqueue(1000 + i)
	}
	assertion.Equal(150, queue.Len())
	{
		expectedQueue := make([]int, 0, 150)
		for i := 0; i < 50; i++ {
			expectedQueue = append(expectedQueue, 50+i)
		}
		for i := 0; i < 100; i++ {
			expectedQueue = append(expectedQueue, 1000+i)
		}
		assertion.Equal(expectedQueue, queue.CopyAsArray())
	}
	// Pop everything
	for i := 0; i < 50; i++ {
		assertion.Equal(50+i, queue.Dequeue())
	}
	for i := 0; i < 100; i++ {
		assertion.Equal(1000+i, queue.Dequeue())
	}
	assertion.Equal(0, queue.Len())
	assertion.Panics(func() { queue.Dequeue() })
}

func TestQueueExists(t *testing.T) {
	queue := NewQueue[int]()
	for i := 0; i < 100; i++ {
		queue.Enqueue(i)
	}
	for i := 0; i < 100; i++ {
		assert.Equal(t, i, queue.Exists(i))
	}
	for i := 0; i < 100; i++ {
		assert.Equal(t, -1, queue.Exists(i+100))
	}
}
