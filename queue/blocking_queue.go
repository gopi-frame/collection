package queue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gopi-frame/contract"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/future"
	"github.com/gopi-frame/utils/catch"
)

// NewBlockingQueue new blocking queue
func NewBlockingQueue[E any](cap int64) *BlockingQueue[E] {
	queue := new(BlockingQueue[E])
	queue.items = []E{}
	queue.cap = cap
	queue.lock = new(sync.RWMutex)
	queue.takeLock = sync.NewCond(queue.lock)
	queue.putLock = sync.NewCond(queue.lock)
	return queue
}

// BlockingQueue blocking queue
type BlockingQueue[E any] struct {
	items    []E
	size     int64
	cap      int64
	takeLock *sync.Cond
	putLock  *sync.Cond
	lock     *sync.RWMutex
}

// Count returns the size of queue
func (q *BlockingQueue[E]) Count() int64 {
	if q.lock.TryRLock() {
		defer q.lock.RUnlock()
	}
	return q.size
}

// IsEmpty returns whether the queue is empty
func (q *BlockingQueue[E]) IsEmpty() bool {
	return q.Count() == 0
}

// IsNotEmpty returns whether the queue is not empty
func (q *BlockingQueue[E]) IsNotEmpty() bool {
	return !q.IsEmpty()
}

// Clear clears the queue
func (q *BlockingQueue[E]) Clear() {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	q.items = nil
	q.size = 0
}

// Peek returns the first element of the queue
func (q *BlockingQueue[E]) Peek() (E, bool) {
	if q.lock.TryRLock() {
		defer q.lock.RUnlock()
	}
	if q.size == 0 {
		return *new(E), false
	}
	return q.items[0], true
}

// TryEnqueue enqueues a new element into the queue, it will return false if the size is up to the capacity
func (q *BlockingQueue[E]) TryEnqueue(value E) bool {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	if q.cap == q.size {
		return false
	}
	q.items = append(q.items, value)
	q.size++
	q.takeLock.Broadcast()
	return true
}

// TryDequeue dequeues the first element of the queue and returns it.
// The empty value of the element type and false will be returned when the queue is empty
func (q *BlockingQueue[E]) TryDequeue() (E, bool) {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	if q.size == 0 {
		return *new(E), false
	}
	value := q.items[0]
	q.items = q.items[1:]
	q.size--
	q.putLock.Broadcast()
	return value, true
}

// Enqueue enqueues a new element into the queue, it will block if the size is up to capacity
func (q *BlockingQueue[E]) Enqueue(value E) bool {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	for q.cap == q.size {
		q.putLock.Wait()
	}
	q.items = append(q.items, value)
	q.size++
	q.takeLock.Broadcast()
	return true
}

// Dequeue dequeues the first element of queue, it will block if the queue is empty
func (q *BlockingQueue[E]) Dequeue() (E, bool) {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	for q.size == 0 {
		q.takeLock.Wait()
	}
	value := q.items[0]
	q.items = q.items[1:]
	q.size--
	q.putLock.Broadcast()
	return value, true
}

// EnqueueTimeout enqueues element into the queue.
// It will block when the size of queue is up to capacity.
// It will return true if the element is successfully enqueued or false when time is out
func (q *BlockingQueue[E]) EnqueueTimeout(value E, duration time.Duration) bool {
	var ok bool
	catch.Try(func() {
		done := make(chan struct{})
		ok = future.Timeout(func() bool {
			future.Void(func() {
				if q.lock.TryLock() {
					defer q.lock.Unlock()
				}
				for q.cap == q.size {
					q.putLock.Wait()
				}
				done <- struct{}{}
			})
			<-done
			q.items = append(q.items, value)
			q.size++
			q.takeLock.Broadcast()
			return true
		}, duration).Complete(func() {
			close(done)
		}).Await()
	}).Catch(new(exception.TimeoutException), func(err error) {
		ok = false
	}).Run()
	return ok
}

// DequeueTimeout removes the first element and returns it.
// It will block when the queue is empty.
// It will return zero value and false when time is out
func (q *BlockingQueue[E]) DequeueTimeout(duration time.Duration) (E, bool) {
	var value E
	var ok bool
	catch.Try(func() {
		done := make(chan struct{})
		ok = future.Timeout(func() bool {
			future.Void(func() {
				if q.lock.TryLock() {
					defer q.lock.Unlock()
				}
				for q.size == 0 {
					q.takeLock.Wait()
				}
				done <- struct{}{}
			})
			<-done
			value = q.items[0]
			q.items = q.items[1:]
			q.size--
			q.putLock.Broadcast()
			return true
		}, duration).Complete(func() {
			close(done)
		}).Await()
	}).Catch(new(exception.TimeoutException), func(err error) {
	}).Run()
	return value, ok
}

// Remove removes the specific element
func (q *BlockingQueue[E]) Remove(value E) {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	var items []E
	for _, item := range q.items {
		if !reflect.DeepEqual(item, value) {
			items = append(items, item)
		}
	}
	q.items = items
	q.size = int64(len(items))
}

// RemoveWhere removes elements which matches the callback
func (q *BlockingQueue[E]) RemoveWhere(callback func(E) bool) {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	var items []E
	for _, item := range q.items {
		if !callback(item) {
			items = append(items, item)
		}
	}
	q.items = items
	q.size = int64(len(items))
}

// ToArray converts to array
func (q *BlockingQueue[E]) ToArray() []E {
	if q.lock.TryRLock() {
		defer q.lock.RUnlock()
	}
	return q.items
}

// ToJSON converts to json
func (q *BlockingQueue[E]) ToJSON() ([]byte, error) {
	return json.Marshal(q.ToArray())
}

// MarshalJSON implements [json.Marshaller]
func (q *BlockingQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (q *BlockingQueue[E]) UnmarshalJSON(data []byte) error {
	if q.lock.TryLock() {
		defer q.lock.Unlock()
	}
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	for _, value := range values {
		for q.size == q.cap {
			q.putLock.Wait()
		}
		q.items = append(q.items, value)
		q.size++
		q.takeLock.Broadcast()
	}
	return nil
}

// String converts to string
func (q *BlockingQueue[E]) String() string {
	if q.lock.TryRLock() {
		defer q.lock.RUnlock()
	}
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("BlockingQueue[%T](len=%d)", *new(E), q.size))
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
	if q.size > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
