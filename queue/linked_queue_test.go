package queue

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedQueue_Count(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	assert.Equal(t, int64(3), queue.Count())
}

func TestLinkedQueue_IsEmpty(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	assert.False(t, queue.IsEmpty())
}

func TestLinkedQueue_IsNotEmpty(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	assert.True(t, queue.IsNotEmpty())
}

func TestLinkedQueue_Clear(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestLinkedQueue_Peek(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	v, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.Equal(t, int64(3), queue.Count())
}

func TestLinkedQueue_Enqueue(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	ok := queue.Enqueue(4)
	assert.True(t, ok)
	assert.Equal(t, int64(4), queue.Count())
	assert.EqualValues(t, []int{1, 2, 3, 4}, queue.ToArray())
}

func TestLinkedQueue_Dequeue(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	v, ok := queue.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
	assert.EqualValues(t, []int{2, 3}, queue.ToArray())
}

func TestLinkedQueue_ToJSON(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	jsonBytes, err := queue.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestLinkedQueue_MarshalJSON(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3)
	jsonBytes, err := json.Marshal(queue)
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestLinkedQueue_UnmarshalJSON(t *testing.T) {
	queue := NewLinkedQueue[int]()
	err := json.Unmarshal([]byte(`[1,2,3]`), queue)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{1, 2, 3}, queue.ToArray())
}

func TestLinkedQueue_String(t *testing.T) {
	queue := NewLinkedQueue(1, 2, 3, 4, 5, 6, 7)
	str := queue.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`LinkedQueue\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\t(\.){3}\n\}`, queue.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}
