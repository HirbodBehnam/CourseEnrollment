package util

import "container/list"

// Queue is a simple linked queue around list.List
type Queue[T comparable] struct {
	l *list.List
}

// NewQueue will create a new queue
func NewQueue[T comparable]() Queue[T] {
	return Queue[T]{list.New()}
}

// Len returns the number of elements in queue
func (q Queue[T]) Len() int {
	return q.l.Len()
}

// Enqueue will push an item into queue
func (q Queue[T]) Enqueue(item T) {
	q.l.PushBack(item)
}

// Dequeue will remove the first of queue and return it
func (q Queue[T]) Dequeue() T {
	front := q.l.Front()
	if front == nil {
		panic("dequeue called on empty queue")
	}
	return q.l.Remove(front).(T)
}

// Exists checks if an item is in the queue. If it exists returns the index of the element in queue.
// Otherwise -1 is returned.
//
// This method runs in O(n) where n is Len
func (q Queue[T]) Exists(item T) int {
	index := 0
	for elem := q.l.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(T) == item {
			return index
		}
		index++
	}
	return -1
}

// Remove will remove a list item based on its value.
// This method will return true if the element is found otherwise false.
//
// This method runs in O(n) where n is Len
func (q Queue[T]) Remove(item T) bool {
	for elem := q.l.Front(); elem != nil; elem = elem.Next() {
		if elem.Value.(T) == item {
			q.l.Remove(elem)
			return true
		}
	}
	return false
}

// CopyAsArray will copy the content of the queue in an array and returns it
func (q Queue[T]) CopyAsArray() []T {
	result := make([]T, 0, q.Len())
	for elem := q.l.Front(); elem != nil; elem = elem.Next() {
		result = append(result, elem.Value.(T))
	}
	return result
}
