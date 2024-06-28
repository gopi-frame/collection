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
	"github.com/gopi-frame/utils/catch"
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

func (q *LinkedBlockingQueue[E]) Count() int64 {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.Count()
}

func (q *LinkedBlockingQueue[E]) IsEmpty() bool {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.IsEmpty()
}

func (q *LinkedBlockingQueue[E]) IsNotEmpty() bool {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.IsNotEmpty()
}

func (q *LinkedBlockingQueue[E]) Clear() {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Clear()
}

func (q *LinkedBlockingQueue[E]) Peek() (E, bool) {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	if q.items.IsEmpty() {
		return *new(E), false
	}
	return q.items.First()
}

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

func (q *LinkedBlockingQueue[E]) Remove(value E) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Remove(value)
}

func (q *LinkedBlockingQueue[E]) RemoveWhere(callback func(E) bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.RemoveWhere(callback)
}

func (q *LinkedBlockingQueue[E]) ToArray() []E {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.ToArray()
}

func (q *LinkedBlockingQueue[E]) ToJSON() ([]byte, error) {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return q.items.MarshalJSON()
}

func (q *LinkedBlockingQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

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
