package queue

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type _comparator struct{}

func (c _comparator) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}

func TestPriorityQueue_Count(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	assert.Equal(t, int64(3), queue.Count())
}

func TestPriorityQueue_IsEmpty(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	assert.False(t, queue.IsEmpty())
}

func TestPriorityQueue_IsNotEmpty(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	assert.True(t, queue.IsNotEmpty())
}

func TestPriorityQueue_Clear(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestPriorityQueue_Peek(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	v, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, int64(3), queue.Count())
}

func TestPriorityQueue_Enqueue(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	ok := queue.Enqueue(4)
	assert.True(t, ok)
	assert.Equal(t, int64(4), queue.Count())
	assert.EqualValues(t, []int{1, 2, 3, 4}, queue.ToArray())
}

func TestPriorityQueue_Dequeue(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	v, ok := queue.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.EqualValues(t, []int{2, 3}, queue.ToArray())
}

func TestPriorityQueue_ToJSON(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	jsonBytes, err := queue.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestPriorityQueue_MarshalJSON(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3)
	jsonBytes, err := json.Marshal(queue)
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestPriorityQueue_UnmarshalJSON(t *testing.T) {
	queue := NewPriorityQueue(_comparator{})
	err := json.Unmarshal([]byte(`[1,2,3]`), queue)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3}, queue.ToArray())
}

func TestPriorityQueue_String(t *testing.T) {
	queue := NewPriorityQueue(_comparator{}, 1, 2, 3, 4, 5, 6, 7)
	str := queue.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`PriorityQueue\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\t(\.){3}\n\}`, queue.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}
