package list

import (
	"encoding/json"
	"fmt"
	"github.com/gopi-frame/exception"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList_IsNotEmpty(t *testing.T) {
	list := NewLinkedList[int](1)
	assert.True(t, list.IsNotEmpty())
}

func TestLinkedList_Contains(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.True(t, list.Contains(1))
}

func TestLinkedList_Remove(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	list.Remove(1)
	assert.False(t, list.Contains(1))
}

func TestLinkedList_RemoveAt(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	list.RemoveAt(0)
	assert.False(t, list.Contains(1))
}

func TestLinkedList_Clear(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	list.Clear()
	assert.True(t, list.IsEmpty())
}

func TestLinkedList_Get(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.Equal(t, 2, list.Get(1))
	assert.PanicsWithError(t, exception.NewRangeException(0, 2).Error(), func() {
		list.Get(3)
	})
}

func TestLinkedList_Set(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	list.Set(0, 2)
	assert.Equal(t, 2, list.Get(0))
}

func TestLinkedList_First(t *testing.T) {
	list := NewLinkedList[int]()
	value, ok := list.First()
	assert.Equal(t, 0, value)
	assert.False(t, ok)
	list.Push(1, 2, 3)
	value, ok = list.First()
	assert.Equal(t, 1, value)
	assert.True(t, ok)
}

func TestLinkedList_FirstOr(t *testing.T) {
	list := NewLinkedList[int]()
	value := list.FirstOr(10)
	assert.Equal(t, 10, value)
	list.Push(0, 1, 2)
	value = list.FirstOr(10)
	assert.Equal(t, 0, value)
}

func TestLinkedList_FirstWhere(t *testing.T) {
	list := NewLinkedList[int](1, 2, 3)
	value, ok := list.FirstWhere(func(item int) bool {
		return item == 3
	})
	assert.Equal(t, 3, value)
	assert.True(t, ok)
	value, ok = list.FirstWhere(func(item int) bool {
		return item > 10
	})
	assert.Equal(t, 0, value)
	assert.False(t, ok)
}

func TestLinkedList_FirstWhereOr(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	value := list.FirstWhereOr(func(item int) bool {
		return item > 10
	}, 10)
	assert.Equal(t, 10, value)
	value = list.FirstWhereOr(func(item int) bool {
		return item >= 2
	}, 10)
	assert.Equal(t, 2, value)
}

func TestLinkedList_Last(t *testing.T) {
	list := NewLinkedList[int]()
	value, ok := list.Last()
	assert.Equal(t, 0, value)
	assert.False(t, ok)
	list.Push(1, 2, 3)
	value, ok = list.Last()
	assert.Equal(t, 3, value)
	assert.True(t, ok)
}

func TestLinkedList_LastOr(t *testing.T) {
	list := NewLinkedList[int]()
	value := list.LastOr(1)
	assert.Equal(t, 1, value)
	list.Push(1, 2, 3)
	value = list.LastOr(2)
	assert.Equal(t, 3, value)
}

func TestLinkedList_LastWhere(t *testing.T) {
	list := NewLinkedList(1, 2, 3, 4)
	value, ok := list.LastWhere(func(item int) bool {
		return item == 3
	})
	assert.Equal(t, 3, value)
	assert.True(t, ok)
	value, ok = list.LastWhere(func(item int) bool {
		return item > 10
	})
	assert.Equal(t, 0, value)
	assert.False(t, ok)
}

func TestLinkedList_LastWhereOr(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	value := list.LastWhereOr(func(item int) bool {
		return item == 4
	}, 10)
	assert.Equal(t, 10, value)
	value = list.LastWhereOr(func(item int) bool {
		return item >= 2
	}, 10)
	assert.Equal(t, 3, value)
}

func TestLinkedList_Pop(t *testing.T) {
	list := NewLinkedList[int]()
	value, ok := list.Pop()
	assert.Equal(t, 0, value)
	assert.False(t, ok)
	list.Push(1, 2, 3)
	value, ok = list.Pop()
	assert.Equal(t, 3, value)
	assert.True(t, ok)
	assert.EqualValues(t, 2, list.Count())
	assert.Equal(t, 2, list.Get(1))
}

func TestLinkedList_Shift(t *testing.T) {
	list := NewLinkedList[int]()
	value, ok := list.Shift()
	assert.Equal(t, 0, value)
	assert.False(t, ok)

	list.Push(1, 2, 3)
	value, ok = list.Shift()
	assert.Equal(t, 1, value)
	assert.True(t, ok)
	assert.EqualValues(t, 2, list.Count())
	assert.Equal(t, 2, list.Get(0))
}

func TestLinkedList_Unshift(t *testing.T) {
	list := NewLinkedList[int](1, 2, 3)
	list.Unshift(0)
	assert.Equal(t, 0, list.Get(0))
	assert.EqualValues(t, 4, list.Count())
}

func TestLinkedList_IndexOf(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.Equal(t, 1, list.IndexOf(2))
}

func TestLinkedList_IndexOfWhere(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.Equal(t, 2, list.IndexOfWhere(func(item int) bool { return item == 3 }))
	assert.Equal(t, -1, list.IndexOfWhere(func(item int) bool { return item > 10 }))
}

func TestLinkedList_Sub(t *testing.T) {
	list := NewLinkedList(1, 2, 3, 4, 5)
	subList := list.Sub(1, 3)
	assert.Equal(t, []int{2, 3}, subList.ToArray())
}

func TestLinkedList_Where(t *testing.T) {
	list := NewLinkedList(1, 2, 3, 4, 5)
	assert.Equal(t, []int{4, 5}, list.Where(func(item int) bool {
		return item > 3
	}).ToArray())
}

func TestLinkedList_Compact(t *testing.T) {
	t.Run("size gte 2", func(t *testing.T) {
		list := NewLinkedList(1, 1, 1, 2, 3, 1, 1)
		list.Compact(nil)
		assert.Equal(t, []int{1, 2, 3, 1}, list.ToArray())
	})

	t.Run("size lt 2", func(t *testing.T) {
		list := NewLinkedList(1)
		list.Compact(nil)
		assert.Equal(t, []int{1}, list.ToArray())
	})
}

func TestLinkedList_Min(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.Equal(t, 1, list.Min(func(a, b int) int {
		if a < b {
			return -1
		}
		return 0
	}))
}

func TestLinkedList_Max(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.Equal(t, 3, list.Max(func(a, b int) int {
		if a < b {
			return -1
		}
		return 1
	}))
}

func TestLinkedList_Sort(t *testing.T) {
	list := NewLinkedList(0, 3, 1, 2)
	list.Sort(func(a, b int) int {
		if a == b {
			return 0
		} else if a < b {
			return -1
		}
		return 1
	})
	assert.Equal(t, []int{0, 1, 2, 3}, list.ToArray())
}

func TestLinkedList_Chunk(t *testing.T) {
	list := NewLinkedList(1, 2, 3, 4)
	chunks := list.Chunk(2)
	assert.EqualValues(t, 2, chunks.Count())
	assert.Equal(t, []any{1, 2}, chunks.Get(0).ToArray())
	assert.Equal(t, []any{3, 4}, chunks.Get(1).ToArray())
}

func TestLinkedList_Each(t *testing.T) {
	list := NewLinkedList(1, 2, 3, 4)
	items := []int{}
	list.Each(func(index, value int) bool {
		items = append(items, value)
		return value < 3
	})
	assert.Equal(t, []int{1, 2, 3}, items)
}

func TestLinkedList_Reverse(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	list.Reverse()
	assert.Equal(t, []int{3, 2, 1}, list.ToArray())
}

func TestLinkedList_Clone(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, list.Clone().ToArray())
}

func TestLinkedList_String(t *testing.T) {
	list := NewLinkedList(1, 2, 3, 4, 5, 6)
	str := list.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`LinkedList\[int\]\(len=%d\)\{\s(\t\d+,\n){5}\t(\.){3}\n\}`, list.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}

func TestLinkedList_ToJSON(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	jsonBytes, err := list.ToJSON()
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
	assert.Nil(t, err)
}

func TestLinkedList_MarshalJSON(t *testing.T) {
	list := NewLinkedList(1, 2, 3)
	jsonBytes, err := json.Marshal(list)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
	assert.Nil(t, err)
}

func TestLinkedList_UnmarshalJSON(t *testing.T) {
	list := NewLinkedList[int]()
	err := json.Unmarshal([]byte(`[1,2,3]`), list)
	assert.Equal(t, []int{1, 2, 3}, list.ToArray())
	assert.Nil(t, err)
}
