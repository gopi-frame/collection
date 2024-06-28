package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type _delay struct {
	value int
	until time.Time
}

func (d _delay) Until() time.Time {
	return d.until
}

func (d _delay) Value() int {
	return d.value
}

func TestDelayedQueue_Count(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(_delay{i, time.Now().Add((5 - 1) * time.Second)})
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_IsEmpty(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	assert.True(t, queue.IsEmpty())
}

func TestDelayedQueue_IsNotEmpty(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(_delay{i, time.Now().Add((5 - 1) * time.Second)})
	}
	assert.True(t, queue.IsNotEmpty())
}

func TestDelayedQueue_Clear(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(_delay{i, time.Now().Add((5 - 1) * time.Second)})
	}
	assert.True(t, queue.IsNotEmpty())
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestDelayedQueue_Peek(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		item := _delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.Enqueue(item)
	}
	value, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 4, value.value)
}

func TestDelayedQueue_TryEnqueue(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		item := _delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.Enqueue(item)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_TryDequeue(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	_, ok := queue.TryDequeue()
	assert.False(t, ok)
	queue.Enqueue(_delay{value: 1, until: time.Now().Add(time.Second)})
	_, ok = queue.TryDequeue()
	assert.False(t, ok)
	time.Sleep(time.Second)
	_, ok = queue.TryDequeue()
	assert.True(t, ok)
}

func TestDelayedQueue_Enqueue(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		item := _delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.Enqueue(item)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_Dequeue(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	queue.Enqueue(_delay{value: 1, until: time.Now().Add(time.Second)})
	start := time.Now()
	_, ok := queue.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, time.Second, time.Second*(time.Duration(time.Since(start).Seconds())))
}

func TestDelayedQueue_EnqueueTimeout(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	for i := 0; i < 5; i++ {
		item := _delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.EnqueueTimeout(item, time.Second)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_DequeueTimeout(t *testing.T) {
	queue := NewDelayedQueue[_delay]()
	_, ok := queue.DequeueTimeout(time.Second)
	assert.False(t, ok)
	go func() {
		queue.Enqueue(_delay{value: 1, until: time.Now().Add(1 * time.Second)})
	}()
	v, ok := queue.DequeueTimeout(2 * time.Second)
	assert.True(t, ok)
	assert.Equal(t, 1, v.Value())
}
