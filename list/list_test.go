package list

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_IsNotEmpty(t *testing.T) {
	list := NewList(1)
	assert.True(t, list.IsNotEmpty())
}

func TestList_Contains(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.True(t, list.Contains(1))
}

func TestList_Remove(t *testing.T) {
	list := NewList(1, 2, 3)
	list.Remove(1)
	assert.False(t, list.Contains(1))
}

func TestList_RemoveAt(t *testing.T) {
	list := NewList(1, 2, 3)
	list.RemoveAt(0)
	assert.False(t, list.Contains(1))
}

func TestList_Clear(t *testing.T) {
	list := NewList(1, 2, 3)
	list.Clear()
	assert.True(t, list.IsEmpty())
}

func TestList_Get(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.Equal(t, 2, list.Get(1))
}

func TestList_Set(t *testing.T) {
	list := NewList(1, 2, 3)
	list.Set(0, 2)
	assert.Equal(t, 2, list.Get(0))
}

func TestList_First(t *testing.T) {
	list := NewList(1, 2, 3)
	value, ok := list.First()
	assert.Equal(t, 1, value)
	assert.True(t, ok)
}

func TestList_FirstOr(t *testing.T) {
	list := NewList[int]()
	value := list.FirstOr(10)
	assert.Equal(t, 10, value)
}

func TestList_FirstWhere(t *testing.T) {
	list := NewList[int](1, 2, 3)
	value, ok := list.FirstWhere(func(item int) bool {
		return item == 3
	})
	assert.Equal(t, 3, value)
	assert.True(t, ok)
}

func TestList_FirstWhereOr(t *testing.T) {
	list := NewList[int]()
	value := list.FirstWhereOr(func(item int) bool {
		return item == 3
	}, 3)
	assert.Equal(t, 3, value)
}

func TestList_Last(t *testing.T) {
	list := NewList(1, 2, 3)
	value, ok := list.Last()
	assert.Equal(t, 3, value)
	assert.True(t, ok)
}

func TestList_LastOr(t *testing.T) {
	list := NewList[int]()
	value := list.LastOr(1)
	assert.Equal(t, 1, value)
}

func TestList_LastWhere(t *testing.T) {
	list := NewList(1, 2, 3, 4)
	value, ok := list.LastWhere(func(item int) bool {
		return item == 3
	})
	assert.Equal(t, 3, value)
	assert.True(t, ok)
}

func TestList_LastWhereOr(t *testing.T) {
	list := NewList(1, 2, 3)
	value := list.LastWhereOr(func(item int) bool {
		return item == 4
	}, 10)
	assert.Equal(t, 10, value)
}

func TestList_Pop(t *testing.T) {
	list := NewList(1, 2, 3)
	value, ok := list.Pop()
	assert.Equal(t, 3, value)
	assert.True(t, ok)
	assert.EqualValues(t, 2, list.Count())
	assert.Equal(t, 2, list.Get(1))
}

func TestList_Shift(t *testing.T) {
	list := NewList(1, 2, 3)
	value, ok := list.Shift()
	assert.Equal(t, 1, value)
	assert.True(t, ok)
	assert.EqualValues(t, 2, list.Count())
	assert.Equal(t, 2, list.Get(0))
}

func TestList_Unshift(t *testing.T) {
	list := NewList(1, 2, 3)
	list.Unshift(0)
	assert.Equal(t, 0, list.Get(0))
	assert.EqualValues(t, 4, list.Count())
}

func TestList_IndexOf(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.Equal(t, 1, list.IndexOf(2))
}

func TestList_IndexOfWhere(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.Equal(t, 2, list.IndexOfWhere(func(item int) bool { return item == 3 }))
}

func TestList_Sub(t *testing.T) {
	list := NewList(1, 2, 3, 4, 5)
	subList := list.Sub(1, 3)
	assert.Equal(t, []int{2, 3}, subList.ToArray())
}

func TestList_Where(t *testing.T) {
	list := NewList(1, 2, 3, 4, 5)
	assert.Equal(t, []int{4, 5}, list.Where(func(item int) bool {
		return item > 3
	}).ToArray())
}

func TestList_Compact(t *testing.T) {
	list := NewList(1, 1, 1, 2, 3, 1, 1)
	list.Compact(nil)
	assert.Equal(t, []int{1, 2, 3, 1}, list.ToArray())
}

func TestList_Min(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.Equal(t, 1, list.Min(func(a, b int) int {
		if a < b {
			return -1
		}
		return 0
	}))
}

func TestList_Max(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.Equal(t, 3, list.Max(func(a, b int) int {
		if a < b {
			return -1
		}
		return 1
	}))
}

func TestList_Sort(t *testing.T) {
	list := NewList(3, 1, 2)
	list.Sort(func(a, b int) int {
		if a == b {
			return 0
		} else if a < b {
			return -1
		}
		return 1
	})
	assert.Equal(t, []int{1, 2, 3}, list.ToArray())
}

func TestList_Chunk(t *testing.T) {
	list := NewList(1, 2, 3, 4)
	chunks := list.Chunk(2)
	assert.EqualValues(t, 2, chunks.Count())
	assert.Equal(t, []any{1, 2}, chunks.Get(0).ToArray())
	assert.Equal(t, []any{3, 4}, chunks.Get(1).ToArray())
}

func TestList_Each(t *testing.T) {
	list := NewList(1, 2, 3, 4)
	items := []int{}
	list.Each(func(index, value int) bool {
		items = append(items, value)
		return value < 3
	})
	assert.Equal(t, []int{1, 2, 3}, items)
}

func TestList_Reverse(t *testing.T) {
	list := NewList(1, 2, 3)
	list.Reverse()
	assert.Equal(t, []int{3, 2, 1}, list.ToArray())
}

func TestList_Clone(t *testing.T) {
	list := NewList(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, list.Clone().ToArray())
}

func TestList_String(t *testing.T) {
	list := NewList(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	str := list.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`List\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\t(\.){3}\n\}`, list.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}

func TestList_ToJSON(t *testing.T) {
	list := NewList(1, 2, 3)
	jsonBytes, err := list.ToJSON()
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
	assert.Nil(t, err)
}

func TestList_MarshalJSON(t *testing.T) {
	list := NewList(1, 2, 3)
	jsonBytes, err := json.Marshal(list)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
	assert.Nil(t, err)
}

func TestList_UnmarshalJSON(t *testing.T) {
	list := NewList[int]()
	err := json.Unmarshal([]byte(`[1,2,3]`), list)
	assert.Equal(t, []int{1, 2, 3}, list.ToArray())
	assert.Nil(t, err)
}
