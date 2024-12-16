package queue

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gopi-frame/contract"
)

// NewPriorityBlockingQueue new priority blocking queue
func NewPriorityBlockingQueue[E any](comparator contract.Comparator[E], cap int64) *PriorityBlockingQueue[E] {
	queue := new(PriorityBlockingQueue[E])
	queue.items = NewPriorityQueue(comparator)
	queue.takeLock = sync.NewCond(queue.items)
	queue.putLock = sync.NewCond(queue.items)
	queue.cap = cap
	return queue
}

// PriorityBlockingQueue priority blocking queue
type PriorityBlockingQueue[E any] struct {
	items    *PriorityQueue[E]
	cap      int64
	takeLock *sync.Cond
	putLock  *sync.Cond
}

// Count returns the size of queue
func (q *PriorityBlockingQueue[E]) Count() int64 {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.Count()
}

// IsEmpty returns whether the queue is empty
func (q *PriorityBlockingQueue[E]) IsEmpty() bool {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.IsEmpty()
}

// IsNotEmpty returns whether the queue is not empty
func (q *PriorityBlockingQueue[E]) IsNotEmpty() bool {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.IsNotEmpty()
}

// Clear clears the queue
func (q *PriorityBlockingQueue[E]) Clear() {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Clear()
}

// Peek returns the first element of the queue
func (q *PriorityBlockingQueue[E]) Peek() (E, bool) {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.Peek()
}

// TryEnqueue enqueues a new element into the queue, it will return false if the size is up to the capacity
func (q *PriorityBlockingQueue[E]) TryEnqueue(value E) bool {
	q.items.RLock()
	defer q.items.RUnlock()
	if q.cap == q.items.Count() {
		return false
	}
	ok := q.items.Enqueue(value)
	q.takeLock.Broadcast()
	return ok
}

// TryDequeue dequeues the first element of the queue and returns it.
// The empty value of the element type and false will be returned when the queue is empty
func (q *PriorityBlockingQueue[E]) TryDequeue() (E, bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	if q.items.Count() == 0 {
		return *new(E), false
	}
	value, ok := q.items.Dequeue()
	q.putLock.Broadcast()
	return value, ok
}

// Enqueue enqueues a new element into the queue, it will block if the size is up to capacity
func (q *PriorityBlockingQueue[E]) Enqueue(value E) bool {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	for q.cap == q.items.Count() {
		q.putLock.Wait()
	}
	ok := q.items.Enqueue(value)
	q.takeLock.Broadcast()
	return ok
}

// Dequeue dequeues the first element of queue, it will block if the queue is empty
func (q *PriorityBlockingQueue[E]) Dequeue() (E, bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	for q.items.IsEmpty() {
		q.takeLock.Wait()
	}
	value, ok := q.items.Dequeue()
	q.putLock.Broadcast()
	return value, ok
}

// EnqueueTimeout enqueues element into the queue.
// It will block when the size of queue is up to capacity.
// It will return true if the element is successfully enqueued or false when time is out
func (q *PriorityBlockingQueue[E]) EnqueueTimeout(value E, duration time.Duration) bool {
	timeout := time.After(duration)
	done := make(chan struct{})
	go func() {
		q.items.Lock()
		defer q.items.Unlock()
		for int64(q.cap) == q.items.Count() {
			q.putLock.Wait()
		}
		close(done)
	}()
	select {
	case <-done:
		ok := q.items.Enqueue(value)
		q.takeLock.Broadcast()
		return ok
	case <-timeout:
		return false
	}
}

// DequeueTimeout removes the first element and returns it.
// It will block when the queue is empty.
// It will return zero value and false when time is out
func (q *PriorityBlockingQueue[E]) DequeueTimeout(duration time.Duration) (value E, ok bool) {
	timeout := time.After(duration)
	done := make(chan struct{})
	go func() {
		q.items.Lock()
		defer q.items.Unlock()
		for q.items.IsEmpty() {
			q.takeLock.Wait()
		}
		close(done)
	}()
	select {
	case <-done:
		value, ok = q.items.Dequeue()
		q.putLock.Broadcast()
		return
	case <-timeout:
		return
	}
}

// Remove removes the specific element
func (q *PriorityBlockingQueue[E]) Remove(value E) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Remove(value)
}

// RemoveWhere removes elements which matches the callback
func (q *PriorityBlockingQueue[E]) RemoveWhere(callback func(E) bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.RemoveWhere(callback)
}

// ToArray converts to array
func (q *PriorityBlockingQueue[E]) ToArray() []E {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.ToArray()
}

// ToJSON converts to json
func (q *PriorityBlockingQueue[E]) ToJSON() ([]byte, error) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	return q.items.ToJSON()
}

// MarshalJSON implements [json.Marshaller]
func (q *PriorityBlockingQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (q *PriorityBlockingQueue[E]) UnmarshalJSON(data []byte) error {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	q.items.Clear()
	for _, value := range values {
		for q.cap == q.items.Count() {
			q.putLock.Wait()
		}
		q.items.Enqueue(value)
		q.takeLock.Broadcast()
	}
	return nil
}

// String converts to string
func (q *PriorityBlockingQueue[E]) String() string {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("PriorityBlockingQueue[%T](len=%d)", *new(E), q.items.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	for index, value := range q.items.items {
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
	if q.items.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
