package queue

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPriorityBlockingQueue_Count(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestPriorityBlockingQueue_IsEmpty(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	assert.True(t, queue.IsEmpty())
}

func TestPriorityBlockingQueue_IsNotEmpty(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.True(t, queue.IsNotEmpty())
}

func TestPriorityBlockingQueue_Clear(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.True(t, queue.IsNotEmpty())
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestPriorityBlockingQueue_Peek(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	value, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 0, value)
	assert.Equal(t, int64(5), queue.Count())
}

func TestPriorityBlockingQueue_TryEnqueue(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		for i := 0; i < 5; i++ {
			queue.Enqueue(i)
		}
		ok := queue.TryEnqueue(6)
		assert.False(t, ok)
	})

	t.Run("empty", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		ok := queue.TryEnqueue(1)
		assert.True(t, ok)
	})
}

func TestPriorityBlockingQueue_TryDequeue(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		_, ok := queue.TryDequeue()
		assert.False(t, ok)
	})

	t.Run("non-empty", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		for i := 0; i < 5; i++ {
			queue.Enqueue(i + 1)
		}
		value, ok := queue.TryDequeue()
		assert.True(t, ok)
		assert.Equal(t, 1, value)
	})
}

func TestPriorityBlockingQueue_Enqueue(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
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

func TestPriorityBlockingQueue_Dequeue(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
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

func TestPriorityBlockingQueue_EnqueueTimeout(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		for i := 0; i < 5; i++ {
			queue.Enqueue(i)
		}
		start := time.Now()
		ok := queue.EnqueueTimeout(6, time.Second)
		assert.False(t, ok)
		assert.Equal(t, time.Second, time.Second*time.Duration(time.Since(start).Seconds()))
	})

	t.Run("non-timeout", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		ok := queue.EnqueueTimeout(6, time.Second)
		assert.True(t, ok)
	})
}

func TestPriorityBlockingQueue_DequeueTimeout(t *testing.T) {
	t.Run("timeout", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		start := time.Now()
		_, ok := queue.DequeueTimeout(time.Second)
		assert.False(t, ok)
		assert.Equal(t, time.Second, time.Second*time.Duration(time.Since(start).Seconds()))
	})

	t.Run("non-timeout", func(t *testing.T) {
		queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
		for i := 0; i < 5; i++ {
			queue.Enqueue(i + 1)
		}
		value, ok := queue.DequeueTimeout(time.Second)
		assert.True(t, ok)
		assert.Equal(t, 1, value)
	})
}

func TestPriorityBlockingQueue_ToArray(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, queue.ToArray())
}

func TestPriorityBlockingQueue_ToJSON(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	jsonBytes, err := queue.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[0,1,2,3,4]`, string(jsonBytes))
}

func TestPriorityBlockingQueue_MarshalJSON(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	jsonBytes, err := json.Marshal(queue)
	assert.Nil(t, err)
	assert.JSONEq(t, `[0,1,2,3,4]`, string(jsonBytes))
}

func TestPriorityBlockingQueue_UnmarshalJSON(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	err := json.Unmarshal([]byte(`[0,1,2,3,4]`), queue)
	assert.Nil(t, err)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, queue.ToArray())
}

func TestPriorityBlockingQueue_String(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	str := queue.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`PriorityBlockingQueue\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\}`, queue.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}

func TestPriorityBlockingQueue_Remove(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	queue.Remove(2)
	assert.Equal(t, int64(4), queue.Count())
	assert.Equal(t, []int{0, 1, 3, 4}, queue.ToArray())
}

func TestPriorityBlockingQueue_RemoveWhere(t *testing.T) {
	queue := NewPriorityBlockingQueue[int](_comparator{}, 5)
	for i := 0; i < 5; i++ {
		queue.Enqueue(i)
	}
	queue.RemoveWhere(func(i int) bool {
		return i%2 == 0
	})
	assert.Equal(t, int64(2), queue.Count())
	assert.Equal(t, []int{1, 3}, queue.ToArray())
}
