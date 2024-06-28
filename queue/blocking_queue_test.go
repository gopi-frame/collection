package queue

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBlockingQueue_Count(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestBlockingQueue_IsEmpty(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	assert.True(t, queue.IsEmpty())
}

func TestBlockingQueue_IsNotEmpty(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.True(t, queue.IsNotEmpty())
}

func TestBlockingQueue_Clear(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.True(t, queue.IsNotEmpty())
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestBlockingQueue_Peek(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	value, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 0, value)
	assert.Equal(t, int64(5), queue.Count())
}

func TestBlockingQueue_TryEnqueue(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	ok := queue.TryEnqueue(6)
	assert.False(t, ok)
}

func TestBlockingQUeue_TryDequeue(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	_, ok := queue.TryDequeue()
	assert.False(t, ok)
}

func TestBlockingQueue_Enqueue(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	go func() {
		start := time.Now()
		ok := queue.Enqueue(6)
		assert.True(t, ok)
		assert.GreaterOrEqual(t, time.Since(start), time.Second)
	}()
	time.Sleep(time.Second)
	queue.Dequeue()
	time.Sleep(time.Second)
}

func TestBlockingQueue_Dequeue(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	go func() {
		start := time.Now()
		v, ok := queue.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
		assert.GreaterOrEqual(t, time.Since(start), time.Second)
	}()
	time.Sleep(time.Second)
	ok := queue.Enqueue(1)
	assert.True(t, ok)
	time.Sleep(time.Second)
}

func TestBlockingQueue_EnqueueTimeout(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	start := time.Now()
	ok := queue.EnqueueTimeout(6, time.Second)
	assert.False(t, ok)
	assert.Equal(t, time.Second, time.Second*time.Duration(time.Since(start).Seconds()))
}

func TestBlockingQueue_DequeueTimeout(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	start := time.Now()
	_, ok := queue.DequeueTimeout(time.Second)
	assert.False(t, ok)
	assert.Equal(t, time.Second, time.Second*time.Duration(time.Since(start).Seconds()))
}

func TestBlockingQueue_ToArray(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, queue.ToArray())
}

func TestBlockingQueue_ToJSON(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	jsonBytes, err := queue.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[0,1,2,3,4]`, string(jsonBytes))
}

func TestBlockingQueue_MarshalJSON(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	jsonBytes, err := json.Marshal(queue)
	assert.Nil(t, err)
	assert.JSONEq(t, `[0,1,2,3,4]`, string(jsonBytes))
}

func TestBlockingQueue_UnmarshalJSON(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	err := json.Unmarshal([]byte(`[0,1,2,3,4]`), queue)
	assert.Nil(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, queue.ToArray())
}

func TestBlockingQueue_String(t *testing.T) {
	queue := NewBlockingQueue[int](5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	str := queue.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`BlockingQueue\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\}`, queue.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}
