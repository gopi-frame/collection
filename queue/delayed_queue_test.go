package queue

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type _delay struct {
	value int
	until time.Time
}

func (d *_delay) Until() time.Time {
	return d.until
}

func (d *_delay) Value() int {
	return d.value
}

func (d *_delay) MarshalJSON() ([]byte, error) {
	type jsonObject struct {
		Value int       `json:"value"`
		Until time.Time `json:"until"`
	}
	return json.Marshal(jsonObject{d.value, d.until})
}

func (d *_delay) UnmarshalJSON(data []byte) error {
	type jsonObject struct {
		Value int       `json:"value"`
		Until time.Time `json:"until"`
	}
	var obj jsonObject
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	d.value = obj.Value
	d.until = obj.Until
	return nil
}

func TestDelayedQueue_Count(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{i, time.Now().Add((5 - 1) * time.Second)})
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_IsEmpty(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	assert.True(t, queue.IsEmpty())
}

func TestDelayedQueue_IsNotEmpty(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{i, time.Now().Add((5 - 1) * time.Second)})
	}
	assert.True(t, queue.IsNotEmpty())
}

func TestDelayedQueue_Clear(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{i, time.Now().Add((5 - 1) * time.Second)})
	}
	assert.True(t, queue.IsNotEmpty())
	queue.Clear()
	assert.True(t, queue.IsEmpty())
}

func TestDelayedQueue_Peek(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		item := &_delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.Enqueue(item)
	}
	value, ok := queue.Peek()
	assert.True(t, ok)
	assert.Equal(t, 4, value.value)
}

func TestDelayedQueue_TryEnqueue(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		item := &_delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.Enqueue(item)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_TryDequeue(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	_, ok := queue.TryDequeue()
	assert.False(t, ok)
	queue.Enqueue(&_delay{value: 1, until: time.Now().Add(time.Second)})
	_, ok = queue.TryDequeue()
	assert.False(t, ok)
	time.Sleep(time.Second)
	_, ok = queue.TryDequeue()
	assert.True(t, ok)
}

func TestDelayedQueue_Enqueue(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		item := &_delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.Enqueue(item)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_Dequeue(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	queue.Enqueue(&_delay{value: 1, until: time.Now().Add(time.Second)})
	start := time.Now()
	_, ok := queue.Dequeue()
	assert.True(t, ok)
	assert.Equal(t, time.Second, time.Second*(time.Duration(time.Since(start).Seconds())))
}

func TestDelayedQueue_EnqueueTimeout(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		item := &_delay{i, time.Now().Add(time.Duration(5-i) * time.Second)}
		queue.EnqueueTimeout(item, time.Second)
	}
	assert.Equal(t, int64(5), queue.Count())
}

func TestDelayedQueue_DequeueTimeout(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	_, ok := queue.DequeueTimeout(time.Second)
	assert.False(t, ok)
	go func() {
		queue.Enqueue(&_delay{value: 1, until: time.Now().Add(1 * time.Second)})
	}()
	v, ok := queue.DequeueTimeout(2 * time.Second)
	assert.True(t, ok)
	assert.Equal(t, 1, v.Value())
}

func TestDelayedQueue_Remove(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	now := time.Now()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{
			value: i,
			until: now.Add(time.Duration(5-i) * time.Second),
		})
	}
	queue.Remove(&_delay{
		value: 2,
		until: now.Add(time.Duration(3) * time.Second),
	})
	assert.Equal(t, int64(4), queue.Count())
}

func TestDelayedQueue_RemoveWhere(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{
			value: i,
			until: time.Now().Add(time.Duration(5-i) * time.Second),
		})
	}
	queue.RemoveWhere(func(value *_delay) bool {
		return value.Value()%2 == 1
	})
	assert.Equal(t, int64(3), queue.Count())
}

func TestDelayedQueue_ToArray(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	now := time.Now()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{
			value: i,
			until: now.Add(time.Duration(5-i) * time.Second),
		})
	}
	var expect []map[string]any
	expect = append(expect, map[string]any{"value": 0, "until": now.Add(5 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 1, "until": now.Add(4 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 2, "until": now.Add(3 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 3, "until": now.Add(2 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 4, "until": now.Add(time.Second).Format(time.RFC3339Nano)})

	var actual []map[string]any
	for _, item := range queue.ToArray() {
		actual = append(actual, map[string]any{"value": item.Value(), "until": item.Until().Format(time.RFC3339Nano)})
	}

	assert.ElementsMatch(t, expect, actual)
}

func TestDelayedQueue_MarshalJSON(t *testing.T) {
	queue := NewDelayedQueue[*_delay]()
	now := time.Now()
	for i := 0; i < 5; i++ {
		queue.Enqueue(&_delay{
			value: i,
			until: now.Add(time.Duration(5-i) * time.Second),
		})
	}
	jsonBytes, err := queue.ToJSON()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	var actual []struct {
		Value int       `json:"value"`
		Until time.Time `json:"until"`
	}
	err = json.Unmarshal(jsonBytes, &actual)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	until1, err := now.Add(5 * time.Second).MarshalJSON()
	until2, err := now.Add(4 * time.Second).MarshalJSON()
	until3, err := now.Add(3 * time.Second).MarshalJSON()
	until4, err := now.Add(2 * time.Second).MarshalJSON()
	until5, err := now.Add(1 * time.Second).MarshalJSON()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	expectedJSONStr := fmt.Sprintf(
		`[
			{"value": 0, "until": %s},
			{"value": 1, "until": %s},
			{"value": 2, "until": %s},
			{"value": 3, "until": %s},
			{"value": 4, "until": %s}
		]`,
		until1, until2, until3, until4, until5,
	)
	var expect []struct {
		Value int       `json:"value"`
		Until time.Time `json:"until"`
	}
	err = json.Unmarshal([]byte(expectedJSONStr), &expect)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	assert.ElementsMatch(t, expect, actual)
}

func TestDelayedQueue_UnmarshalJSON(t *testing.T) {
	now := time.Now()
	until1, err := now.Add(5 * time.Second).MarshalJSON()
	until2, err := now.Add(4 * time.Second).MarshalJSON()
	until3, err := now.Add(3 * time.Second).MarshalJSON()
	until4, err := now.Add(2 * time.Second).MarshalJSON()
	until5, err := now.Add(1 * time.Second).MarshalJSON()
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	jsonBytes := fmt.Sprintf(
		`[
			{"value": 0, "until": %s},
			{"value": 1, "until": %s},
			{"value": 2, "until": %s},
			{"value": 3, "until": %s},
			{"value": 4, "until": %s}
		]`,
		until1, until2, until3, until4, until5,
	)
	var queue = NewDelayedQueue[*_delay]()
	err = json.Unmarshal([]byte(jsonBytes), queue)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	var expect []map[string]any
	expect = append(expect, map[string]any{"value": 0, "until": now.Add(5 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 1, "until": now.Add(4 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 2, "until": now.Add(3 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 3, "until": now.Add(2 * time.Second).Format(time.RFC3339Nano)})
	expect = append(expect, map[string]any{"value": 4, "until": now.Add(1 * time.Second).Format(time.RFC3339Nano)})

	var actual []map[string]any
	for _, item := range queue.ToArray() {
		actual = append(actual, map[string]any{"value": item.Value(), "until": item.Until().Format(time.RFC3339Nano)})
	}

	assert.ElementsMatch(t, expect, actual)
}
