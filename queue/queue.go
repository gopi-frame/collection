package queue

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/contract"
)

// NewQueue new queue
func NewQueue[E any](values ...E) *Queue[E] {
	queue := new(Queue[E])
	queue.items = list.NewList(values...)
	return queue
}

// Queue array queue
type Queue[E any] struct {
	items *list.List[E]
}

// Lock locks the queue
func (q *Queue[E]) Lock() {
	q.items.Lock()
}

// Unlock unlocks the queue
func (q *Queue[E]) Unlock() {
	q.items.Unlock()
}

// TryLock tries to lock the queue
func (q *Queue[E]) TryLock() bool {
	return q.items.TryLock()
}

// RLock locks the read lock for the queue
func (q *Queue[E]) RLock() {
	q.items.RLock()
}

// TryRLock tries to lock the read lock for the queue
func (q *Queue[E]) TryRLock() bool {
	return q.items.TryRLock()
}

// RUnlock unlocks the read lock for the queue
func (q *Queue[E]) RUnlock() {
	q.items.RUnlock()
}

// Count returns the size of queue
func (q *Queue[E]) Count() int64 {
	return q.items.Count()
}

// IsEmpty returns whether the queue is empty
func (q *Queue[E]) IsEmpty() bool {
	return q.Count() == 0
}

// IsNotEmpty returns whether the queue is not empty
func (q *Queue[E]) IsNotEmpty() bool {
	return !q.IsEmpty()
}

// Clear clears the queue
func (q *Queue[E]) Clear() {
	q.items.Clear()
}

// Peek returns the first element of the queue
func (q *Queue[E]) Peek() (E, bool) {
	return q.items.First()
}

// Enqueue enqueues a new element into the queue, it will block if the size is up to capacity
func (q *Queue[E]) Enqueue(value E) bool {
	q.items.Push(value)
	return true
}

// Dequeue dequeues the first element of queue, it will block if the queue is empty
func (q *Queue[E]) Dequeue() (E, bool) {
	return q.items.Shift()
}

// Remove removes the specific element
func (q *Queue[E]) Remove(value E) {
	q.items.Remove(value)
}

// RemoveWhere removes elements which matches the callback
func (q *Queue[E]) RemoveWhere(callback func(value E) bool) {
	q.items.RemoveWhere(callback)
}

// ToArray converts to array
func (q *Queue[E]) ToArray() []E {
	return q.items.ToArray()
}

// ToJSON converts to json
func (q *Queue[E]) ToJSON() ([]byte, error) {
	return q.items.ToJSON()
}

// MarshalJSON implements [json.Marshaller]
func (q *Queue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (q *Queue[E]) UnmarshalJSON(data []byte) error {
	var values []E
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}
	q.items = list.NewList[E](values...)
	return nil
}

// String converts to string
func (q *Queue[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("Queue[%T](len=%d)", *new(E), q.Count()))
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
