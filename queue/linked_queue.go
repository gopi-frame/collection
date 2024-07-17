package queue

import (
	"fmt"
	"strings"

	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/contract"
)

// NewLinkedQueue new linked queue
func NewLinkedQueue[E any](values ...E) *LinkedQueue[E] {
	queue := new(LinkedQueue[E])
	queue.items = list.NewLinkedList(values...)
	return queue
}

// LinkedQueue linked queue
type LinkedQueue[E any] struct {
	items *list.LinkedList[E]
}

// Lock locks the queue
func (q *LinkedQueue[E]) Lock() {
	q.items.Lock()
}

// Unlock unlocks the queue
func (q *LinkedQueue[E]) Unlock() {
	q.items.Unlock()
}

// TryLock tries to lock the queue
func (q *LinkedQueue[E]) TryLock() bool {
	return q.items.TryLock()
}

// RLock locks the read lock for the queue
func (q *LinkedQueue[E]) RLock() {
	q.items.RLock()
}

// RUnlock unlocks the read lock for the queue
func (q *LinkedQueue[E]) RUnlock() {
	q.items.RUnlock()
}

// TryRLock tries to lock the read lock for the queue
func (q *LinkedQueue[E]) TryRLock() bool {
	return q.items.TryRLock()
}

// Count returns the size of queue
func (q *LinkedQueue[E]) Count() int64 {
	return q.items.Count()
}

// IsEmpty returns whether the queue is empty
func (q *LinkedQueue[E]) IsEmpty() bool {
	return q.items.IsEmpty()
}

// IsNotEmpty returns whether the queue is not empty
func (q *LinkedQueue[E]) IsNotEmpty() bool {
	return q.items.IsNotEmpty()
}

// Clear clears the queue
func (q *LinkedQueue[E]) Clear() {
	q.items.Clear()
}

// Peek returns the first element of the queue
func (q *LinkedQueue[E]) Peek() (E, bool) {
	return q.items.First()
}

// Enqueue enqueues a new element into the queue, it will block if the size is up to capacity
func (q *LinkedQueue[E]) Enqueue(value E) bool {
	q.items.Push(value)
	return true
}

// Dequeue dequeues the first element of queue, it will block if the queue is empty
func (q *LinkedQueue[E]) Dequeue() (value E, ok bool) {
	if q.items.IsEmpty() {
		return
	}
	return q.items.Shift()
}

// Remove removes the specific element
func (q *LinkedQueue[E]) Remove(value E) {
	q.items.Remove(value)
}

// RemoveWhere removes elements which matches the callback
func (q *LinkedQueue[E]) RemoveWhere(callback func(value E) bool) {
	q.items.RemoveWhere(callback)
}

// ToArray converts to array
func (q *LinkedQueue[E]) ToArray() []E {
	return q.items.ToArray()
}

// ToJSON converts to json
func (q *LinkedQueue[E]) ToJSON() ([]byte, error) {
	return q.items.MarshalJSON()
}

// MarshalJSON implements [json.Marshaller]
func (q *LinkedQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (q *LinkedQueue[E]) UnmarshalJSON(data []byte) error {
	return q.items.UnmarshalJSON(data)
}

// String converts to string
func (q *LinkedQueue[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedQueue[%T](len=%d)", *new(E), q.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	q.items.Each(func(index int, value E) bool {
		str.WriteByte('\t')
		if v, ok := any(value).(contract.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		return index < 4
	})
	if q.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
