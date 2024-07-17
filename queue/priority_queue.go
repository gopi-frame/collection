package queue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/gopi-frame/contract"
)

// NewPriorityQueue new priority queue
func NewPriorityQueue[E any](comparator contract.Comparator[E], values ...E) *PriorityQueue[E] {
	queue := new(PriorityQueue[E])
	queue.comparator = comparator
	for _, value := range values {
		queue.Enqueue(value)
	}
	return queue
}

// PriorityQueue priority queue
type PriorityQueue[E any] struct {
	sync.RWMutex
	size       int64
	items      []E
	comparator contract.Comparator[E]
}

func (q *PriorityQueue[E]) less(i, j int64) bool {
	return q.comparator.Compare(q.items[i], q.items[j]) < 0
}

func (q *PriorityQueue[E]) swap(i, j int64) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

// Count returns the size of queue
func (q *PriorityQueue[E]) Count() int64 {
	return q.size
}

// IsEmpty returns whether the queue is empty
func (q *PriorityQueue[E]) IsEmpty() bool {
	return q.Count() == 0
}

// IsNotEmpty returns whether the queue is not empty
func (q *PriorityQueue[E]) IsNotEmpty() bool {
	return !q.IsEmpty()
}

// Clear clears the queue
func (q *PriorityQueue[E]) Clear() {
	q.items = make([]E, 0)
	q.size = 0
}

// Peek returns the first element of the queue
func (q *PriorityQueue[E]) Peek() (E, bool) {
	if q.size == 0 {
		return *new(E), false
	}
	return q.items[0], true
}

// Enqueue enqueues a new element into the queue, it will block if the size is up to capacity
func (q *PriorityQueue[E]) Enqueue(value E) bool {
	q.items = append(q.items, value)
	q.size++
	for index := q.size - 1; q.less(index, (index-1)/2); index = (index - 1) / 2 {
		q.swap(index, (index-1)/2)
	}
	return true
}

// Dequeue dequeues the first element of queue, it will block if the queue is empty
func (q *PriorityQueue[E]) Dequeue() (value E, ok bool) {
	if q.size == 0 {
		return *new(E), false
	}
	value = q.items[0]
	ok = true
	q.swap(0, q.size-1)
	q.items = q.items[:q.size-1]
	q.size--
	index := int64(0)
	lastIndex := q.size - 1
	for {
		leftIndex := index*2 + 1
		if leftIndex > lastIndex || leftIndex < 0 {
			break
		}
		swapIndex := leftIndex
		if rightIndex := leftIndex + 1; rightIndex <= lastIndex && q.less(rightIndex, leftIndex) {
			swapIndex = rightIndex
		}
		if !q.less(swapIndex, index) {
			break
		}
		q.swap(swapIndex, index)
		index = swapIndex
	}
	return
}

// Remove removes the specific element
func (q *PriorityQueue[E]) Remove(value E) {
	q.RemoveWhere(func(e E) bool {
		return reflect.DeepEqual(e, value)
	})
}

// RemoveWhere removes elements which matches the callback
func (q *PriorityQueue[E]) RemoveWhere(callback func(E) bool) {
	q.items = slices.DeleteFunc(q.items, callback)
	q.size = int64(len(q.items))
}

// ToArray converts to array
func (q *PriorityQueue[E]) ToArray() []E {
	return q.items
}

// ToJSON converts to json
func (q *PriorityQueue[E]) ToJSON() ([]byte, error) {
	return json.Marshal(q.ToArray())
}

// MarshalJSON implements [json.Marshaller]
func (q *PriorityQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (q *PriorityQueue[E]) UnmarshalJSON(data []byte) error {
	items := []E{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return nil
	}
	q.Clear()
	for _, item := range items {
		q.Enqueue(item)
	}
	return nil
}

// String converts to string
func (q *PriorityQueue[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("PriorityQueue[%T](len=%d)", *new(E), q.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	for index, value := range q.items {
		str.WriteByte('\t')
		if v, ok := any(value).(contract.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		if index >= 4 {
			break
		}
	}
	if q.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
