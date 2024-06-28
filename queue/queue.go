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

func (q *Queue[E]) Lock() {
	q.items.Lock()
}

func (q *Queue[E]) Unlock() {
	q.items.Unlock()
}

func (q *Queue[E]) TryLock() bool {
	return q.items.TryLock()
}

func (q *Queue[E]) RLock() {
	q.items.RLock()
}

func (q *Queue[E]) TryRLock() bool {
	return q.items.TryRLock()
}

func (q *Queue[E]) RUnlock() {
	q.items.RUnlock()
}

func (q *Queue[E]) Count() int64 {
	return q.items.Count()
}

func (q *Queue[E]) IsEmpty() bool {
	return q.Count() == 0
}

func (q *Queue[E]) IsNotEmpty() bool {
	return !q.IsEmpty()
}

func (q *Queue[E]) Clear() {
	q.items.Clear()
}

func (q *Queue[E]) Peek() (E, bool) {
	return q.items.First()
}

func (q *Queue[E]) Enqueue(value E) bool {
	q.items.Push(value)
	return true
}

func (q *Queue[E]) Dequeue() (E, bool) {
	return q.items.Shift()
}

func (q *Queue[E]) Remove(value E) {
	q.items.Remove(value)
}

func (q *Queue[E]) RemoveWhere(callback func(value E) bool) {
	q.items.RemoveWhere(callback)
}

func (q *Queue[E]) ToArray() []E {
	return q.items.ToArray()
}

func (q *Queue[E]) ToJSON() ([]byte, error) {
	return q.items.ToJSON()
}

func (q *Queue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

func (q *Queue[E]) UnmarshalJSON(data []byte) error {
	var values = []E{}
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}
	q.items = list.NewList[E](values...)
	return nil
}

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
