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

func (q *LinkedQueue[E]) Lock() {
	q.items.Lock()
}

func (q *LinkedQueue[E]) Unlock() {
	q.items.Unlock()
}

func (q *LinkedQueue[E]) TryLock() bool {
	return q.items.TryLock()
}

func (q *LinkedQueue[E]) RLock() {
	q.items.RLock()
}

func (q *LinkedQueue[E]) RUnlock() {
	q.items.RUnlock()
}

func (q *LinkedQueue[E]) TryRLock() bool {
	return q.items.TryRLock()
}

func (q *LinkedQueue[E]) Count() int64 {
	return q.items.Count()
}

func (q *LinkedQueue[E]) IsEmpty() bool {
	return q.items.IsEmpty()
}

func (q *LinkedQueue[E]) IsNotEmpty() bool {
	return q.items.IsNotEmpty()
}

func (q *LinkedQueue[E]) Clear() {
	q.items.Clear()
}

func (q *LinkedQueue[E]) Peek() (E, bool) {
	return q.items.First()
}

func (q *LinkedQueue[E]) Enqueue(value E) bool {
	q.items.Push(value)
	return true
}

func (q *LinkedQueue[E]) Dequeue() (value E, ok bool) {
	if q.items.IsEmpty() {
		return
	}
	return q.items.Shift()
}

func (q *LinkedQueue[E]) Remove(value E) {
	q.items.Remove(value)
}

func (q *LinkedQueue[E]) RemoveWhere(callback func(value E) bool) {
	q.items.RemoveWhere(callback)
}

func (q *LinkedQueue[E]) ToArray() []E {
	return q.items.ToArray()
}

func (q *LinkedQueue[E]) ToJSON() ([]byte, error) {
	return q.items.MarshalJSON()
}

func (q *LinkedQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

func (q *LinkedQueue[E]) UnmarshalJSON(data []byte) error {
	return q.items.UnmarshalJSON(data)
}

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
