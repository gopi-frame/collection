package kv

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_IsNotEmpty(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.IsNotEmpty())
}

func TestMap_Get(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	v, ok := m.Get(1)
	assert.Equal(t, 1, v)
	assert.True(t, ok)
}

func TestMap_GetOr(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	v := m.GetOr(10, 1)
	assert.Equal(t, 1, v)
	v = m.GetOr(1, 10)
	assert.Equal(t, 1, v)
}

func TestMap_Remove(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	m.Remove(1)
	v, ok := m.Get(1)
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestMap_Keys(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.ElementsMatch(t, []int{0, 1, 2}, m.Keys())
}

func TestMap_Values(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.ElementsMatch(t, []int{0, 1, 2}, m.Values())
}

func TestMap_Clear(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	m.Clear()
	assert.True(t, m.IsEmpty())
}

func TestMap_ContainsKey(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.ContainsKey(0))
	assert.False(t, m.ContainsKey(-1))
}

func TestMap_Contains(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.Contains(0))
}

func TestMap_ContainsWhere(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.ContainsWhere(func(value int) bool {
		return value == 2
	}))
	assert.False(t, m.ContainsWhere(func(value int) bool {
		return value < 0
	}))
}

func TestMap_Each(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	t.Run("without break", func(t *testing.T) {
		var items []int
		m.Each(func(key, value int) bool {
			if value <= 1 {
				items = append(items, value)
			}
			return true
		})
		assert.ElementsMatch(t, []int{0, 1}, items)
	})

	t.Run("with break", func(t *testing.T) {
		var items []int
		m.Each(func(key, value int) bool {
			return false
		})
		assert.Len(t, items, 0)
	})
}

func TestMap_ToJSON(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	jsonBytes, err := m.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `{"0":0,"1":1,"2":2}`, string(jsonBytes))
}

func TestMap_MarshalJSON(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	jsonBytes, err := json.Marshal(m)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"0":0,"1":1,"2":2}`, string(jsonBytes))
}

func TestMap_UnmarshalJSON(t *testing.T) {
	m := NewMap[int, int]()
	err := json.Unmarshal([]byte(`{"0":0,"1":1,"2":2}`), m)
	assert.Nil(t, err)
	assert.EqualValues(t, map[int]int{
		0: 0, 1: 1, 2: 2,
	}, m.ToMap())
}

func TestMap_String(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	str := m.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`Map\[int, int\]\(len=%d\)\{\n(\t\d+:\s\d+,\n)+\}`, m.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}

func TestMap_Clone(t *testing.T) {
	m := NewMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	m2 := m.Clone()
	assert.EqualValues(t, map[int]int{
		0: 0, 1: 1, 2: 2,
	}, m2.ToMap())
}
