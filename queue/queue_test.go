package queue

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"sync"
	"testing"
)

func TestQueue_Count(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	assert.Equal(t, int64(3), queue.Count())
}

func TestQueue_IsEmpty(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	assert.False(t, queue.IsEmpty())
}

func TestQueue_IsNotEmpty(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	assert.True(t, queue.IsNotEmpty())
}

func TestQueue_Clear(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestQueue_Peek(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	v, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, int64(3), queue.Count())
}

func TestQueue_Enqueue(t *testing.T) {
	t.Run("standalone-coroutine", func(t *testing.T) {
		queue := NewQueue(1, 2, 3)
		ok := queue.Enqueue(4)
		assert.True(t, ok)
		assert.Equal(t, int64(4), queue.Count())
		assert.Equal(t, []int{1, 2, 3, 4}, queue.ToArray())
	})

	t.Run("multi-coroutines", func(t *testing.T) {
		queue := NewQueue[int]()
		var expected []int
		var wg = new(sync.WaitGroup)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			expected = append(expected, i)
			go func(i int) {
				queue.Lock()
				assert.True(t, queue.Enqueue(i))
				queue.Unlock()
				wg.Done()
			}(i)
		}
		wg.Wait()
		assert.ElementsMatch(t, expected, queue.ToArray())
	})
}

func TestQueue_Dequeue(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	v, ok := queue.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.EqualValues(t, []int{2, 3}, queue.ToArray())
}

func TestQueue_ToJSON(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	jsonBytes, err := queue.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestQueue_MarshalJSON(t *testing.T) {
	queue := NewQueue(1, 2, 3)
	jsonBytes, err := json.Marshal(queue)
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestQueue_UnmarshalJSON(t *testing.T) {
	queue := NewQueue[int]()
	err := json.Unmarshal([]byte(`[1,2,3]`), queue)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3}, queue.ToArray())
}

func TestQueue_String(t *testing.T) {
	queue := NewQueue(1, 2, 3, 4, 5, 6, 7)
	str := queue.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`Queue\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\t(\.){3}\n\}`, queue.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}

func TestQueue_Remove(t *testing.T) {
	queue := NewQueue(1, 2, 3, 4, 5)
	queue.Remove(1)
	assert.Equal(t, int64(4), queue.Count())
	assert.Equal(t, []int{2, 3, 4, 5}, queue.ToArray())
}

func TestQueue_RemoveWhere(t *testing.T) {
	queue := NewQueue(1, 2, 3, 4, 5)
	queue.RemoveWhere(func(value int) bool {
		return value%2 == 1
	})
	assert.Equal(t, int64(2), queue.Count())
	assert.Equal(t, []int{2, 4}, queue.ToArray())
}
