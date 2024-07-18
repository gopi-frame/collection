package queue

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/contract"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/future"
	"github.com/gopi-frame/util/catch"
)

// NewLinkedBlockingQueue new linked blocking queue
func NewLinkedBlockingQueue[E any](cap int) *LinkedBlockingQueue[E] {
	queue := new(LinkedBlockingQueue[E])
	queue.items = list.NewLinkedList[E]()
	queue.takeLock = sync.NewCond(queue.items)
	queue.putLock = sync.NewCond(queue.items)
	queue.cap = cap
	return queue
}

// LinkedBlockingQueue linked blocking queue
type LinkedBlockingQueue[E any] struct {
	items    *list.LinkedList[E]
	cap      int
	takeLock *sync.Cond
	putLock  *sync.Cond
}

// Count returns the size of queue
func (q *LinkedBlockingQueue[E]) Count() int64 {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.Count()
}

// IsEmpty returns whether the queue is empty
func (q *LinkedBlockingQueue[E]) IsEmpty() bool {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.IsEmpty()
}

// IsNotEmpty returns whether the queue is not empty
func (q *LinkedBlockingQueue[E]) IsNotEmpty() bool {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.IsNotEmpty()
}

// Clear clears the queue
func (q *LinkedBlockingQueue[E]) Clear() {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Clear()
}

// Peek returns the first element of the queue
func (q *LinkedBlockingQueue[E]) Peek() (E, bool) {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	if q.items.IsEmpty() {
		return *new(E), false
	}
	return q.items.First()
}

// TryEnqueue enqueues a new element into the queue, it will return false if the size is up to the capacity
func (q *LinkedBlockingQueue[E]) TryEnqueue(value E) bool {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	if int64(q.cap) == q.items.Count() {
		return false
	}
	q.items.Push(value)
	q.takeLock.Broadcast()
	return true
}

// TryDequeue dequeues the first element of the queue and returns it.
// The empty value of the element type and false will be returned when the queue is empty
func (q *LinkedBlockingQueue[E]) TryDequeue() (E, bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	if q.items.IsEmpty() {
		return *new(E), false
	}
	value, ok := q.items.Shift()
	q.putLock.Broadcast()
	return value, ok
}

// Enqueue enqueues a new element into the queue, it will block if the size is up to capacity
func (q *LinkedBlockingQueue[E]) Enqueue(value E) bool {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	for int64(q.cap) == q.items.Count() {
		q.putLock.Wait()
	}
	q.items.Push(value)
	q.takeLock.Broadcast()
	return true
}

// Dequeue dequeues the first element of queue, it will block if the queue is empty
func (q *LinkedBlockingQueue[E]) Dequeue() (E, bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	for q.items.IsEmpty() {
		q.takeLock.Wait()
	}
	value, ok := q.items.Shift()
	q.putLock.Broadcast()
	return value, ok
}

// EnqueueTimeout enqueues element into the queue.
// It will block when the size of queue is up to capacity.
// It will return true if the element is successfully enqueued or false when time is out
func (q *LinkedBlockingQueue[E]) EnqueueTimeout(value E, duration time.Duration) bool {
	var ok bool
	catch.Try(func() {
		done := make(chan struct{})
		ok = future.Timeout(func() bool {
			future.Void(func() {
				q.items.Lock()
				defer q.items.Unlock()
				for int64(q.cap) == q.items.Count() {
					q.putLock.Wait()
				}
				done <- struct{}{}
			})
			<-done
			q.items.Push(value)
			q.takeLock.Broadcast()
			return true
		}, duration).Complete(func() {
			close(done)
		}).Await()
	}).Catch(new(exception.TimeoutException), func(err error) {
	}).Run()
	return ok
}

// DequeueTimeout removes the first element and returns it.
// It will block when the queue is empty.
// It will return zero value and false when time is out
func (q *LinkedBlockingQueue[E]) DequeueTimeout(duration time.Duration) (E, bool) {
	var value E
	var ok bool
	catch.Try(func() {
		done := make(chan struct{})
		future.Timeout(func() bool {
			future.Void(func() {
				q.items.Lock()
				defer q.items.Unlock()
				for q.items.IsEmpty() {
					q.takeLock.Wait()
				}
				done <- struct{}{}
			})
			<-done
			value, ok = q.items.Shift()
			q.putLock.Broadcast()
			return ok
		}, duration).Complete(func() {
			close(done)
		}).Await()
	}).Catch(new(exception.TimeoutException), func(err error) {
	}).Run()
	return value, ok
}

// Remove removes the specific element
func (q *LinkedBlockingQueue[E]) Remove(value E) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Remove(value)
}

// RemoveWhere removes elements which matches the callback
func (q *LinkedBlockingQueue[E]) RemoveWhere(callback func(E) bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.RemoveWhere(callback)
}

// ToArray converts to array
func (q *LinkedBlockingQueue[E]) ToArray() []E {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.ToArray()
}

// ToJSON converts to json
func (q *LinkedBlockingQueue[E]) ToJSON() ([]byte, error) {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.MarshalJSON()
}

// MarshalJSON implements [json.Marshaller]
func (q *LinkedBlockingQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (q *LinkedBlockingQueue[E]) UnmarshalJSON(data []byte) error {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	for _, value := range values {
		for q.items.Count() == int64(q.cap) {
			q.putLock.Wait()
		}
		q.items.Push(value)
		q.takeLock.Broadcast()
	}
	return nil
}

// String converts to string
func (q *LinkedBlockingQueue[E]) String() string {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedBlockingQueue[%T](len=%d)", *new(E), q.items.Count()))
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
	if q.items.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
