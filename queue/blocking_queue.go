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

// NewBlockingQueue new blocking queue
func NewBlockingQueue[E any](cap int) *BlockingQueue[E] {
	queue := new(BlockingQueue[E])
	queue.items = list.NewList(make([]E, cap)...)
	queue.cap = cap
	queue.takeLock = sync.NewCond(queue.items)
	queue.putLock = sync.NewCond(queue.items)
	return queue
}

// BlockingQueue blocking queue
type BlockingQueue[E any] struct {
	items        *list.List[E]
	size         int
	cap          int
	enqueueIndex int
	dequeueIndex int
	takeLock     *sync.Cond
	putLock      *sync.Cond
}

func (q *BlockingQueue[E]) moveIndex(index int) int {
	index++
	if index >= q.cap {
		return 0
	}
	return index
}

func (q *BlockingQueue[E]) Count() int64 {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	return int64(q.size)
}

func (q *BlockingQueue[E]) IsEmpty() bool {
	return q.Count() == 0
}

func (q *BlockingQueue[E]) IsNotEmpty() bool {
	return !q.IsEmpty()
}

func (q *BlockingQueue[E]) Clear() {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	q.items.Clear()
	q.size = 0
	q.dequeueIndex = 0
	q.enqueueIndex = 0
}

func (q *BlockingQueue[E]) Peek() (E, bool) {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	if q.size == 0 {
		return *new(E), false
	}
	return q.items.Get(q.dequeueIndex), true
}

func (q *BlockingQueue[E]) TryEnqueue(value E) bool {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	if q.cap == q.size {
		return false
	}
	q.items.Set(q.enqueueIndex, value)
	q.size++
	q.enqueueIndex = q.moveIndex(q.enqueueIndex)
	q.takeLock.Broadcast()
	return true
}

func (q *BlockingQueue[E]) TryDequeue() (E, bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	if q.size == 0 {
		return *new(E), false
	}
	value := q.items.Get(q.dequeueIndex)
	q.items.Set(q.dequeueIndex, *new(E))
	q.dequeueIndex = q.moveIndex(q.dequeueIndex)
	q.size--
	q.putLock.Broadcast()
	return value, true
}

func (q *BlockingQueue[E]) Enqueue(value E) bool {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	for q.cap == q.size {
		q.putLock.Wait()
	}
	q.items.Set(q.enqueueIndex, value)
	q.size++
	q.enqueueIndex = q.moveIndex(q.enqueueIndex)
	q.takeLock.Broadcast()
	return true
}

func (q *BlockingQueue[E]) Dequeue() (E, bool) {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	for q.size == 0 {
		q.takeLock.Wait()
	}
	value := q.items.Get(q.dequeueIndex)
	q.items.Set(q.dequeueIndex, *new(E))
	q.size--
	q.dequeueIndex = q.moveIndex(q.dequeueIndex)
	q.putLock.Broadcast()
	return value, true
}

func (q *BlockingQueue[E]) EnqueueTimeout(value E, duration time.Duration) bool {
	var ok bool
	catch.Try(func() {
		done := make(chan struct{})
		ok = future.Timeout(func() bool {
			future.Void(func() {
				q.items.Lock()
				defer q.items.Unlock()
				for q.cap == q.size {
					q.putLock.Wait()
				}
				done <- struct{}{}
			})
			<-done
			q.items.Set(q.enqueueIndex, value)
			q.size++
			q.enqueueIndex = q.moveIndex(q.enqueueIndex)
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

func (q *BlockingQueue[E]) DequeueTimeout(duration time.Duration) (E, bool) {
	var value E
	var ok bool
	catch.Try(func() {
		done := make(chan struct{})
		ok = future.Timeout(func() bool {
			future.Void(func() {
				q.items.Lock()
				defer q.items.Unlock()
				for q.size == 0 {
					q.takeLock.Wait()
				}
				done <- struct{}{}
			})
			<-done
			value = q.items.Get(q.dequeueIndex)
			q.items.Set(q.dequeueIndex, *new(E))
			q.size--
			q.dequeueIndex = q.moveIndex(q.dequeueIndex)
			q.putLock.Broadcast()
			return true
		}, duration).Complete(func() {
			close(done)
		}).Await()
	}).Catch(new(exception.TimeoutException), func(err error) {
	}).Run()
	return value, ok
}

func (q *BlockingQueue[E]) ToArray() []E {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	values := []E{}
	if q.enqueueIndex > q.dequeueIndex {
		for index := q.dequeueIndex; index < q.enqueueIndex; index++ {
			values = append(values, q.items.Get(index))
		}
	} else {
		for index := q.dequeueIndex; index < q.cap; index++ {
			values = append(values, q.items.Get(index))
		}
		for index := 0; index < q.enqueueIndex; index++ {
			values = append(values, q.items.Get(index))
		}
	}
	return values
}

func (q *BlockingQueue[E]) ToJSON() ([]byte, error) {
	return json.Marshal(q.ToArray())
}

func (q *BlockingQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

func (q *BlockingQueue[E]) UnmarshalJSON(data []byte) error {
	if q.items.TryLock() {
		defer q.items.Unlock()
	}
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	for _, value := range values {
		for q.size == q.cap {
			q.putLock.Wait()
		}
		q.items.Set(q.enqueueIndex, value)
		q.size++
		q.enqueueIndex = q.moveIndex(q.enqueueIndex)
		q.takeLock.Broadcast()
	}
	return nil
}

func (q *BlockingQueue[E]) String() string {
	if q.items.TryRLock() {
		defer q.items.RUnlock()
	}
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("BlockingQueue[%T](len=%d)", *new(E), q.size))
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
	if q.size > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
