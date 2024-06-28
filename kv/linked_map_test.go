package kv

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedMap_IsNotEmpty(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.IsNotEmpty())
}

func TestLinkedMap_Get(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	v, ok := m.Get(1)
	assert.Equal(t, 1, v)
	assert.True(t, ok)
}

func TestLinkedMap_GetOr(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	v := m.GetOr(10, 1)
	assert.Equal(t, 1, v)
}

func TestLinkedMap_Remove(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	m.Remove(1)
	v, ok := m.Get(1)
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestLinkedMap_First(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(2, 2)
	m.Set(0, 0)
	m.Set(1, 1)
	v, ok := m.First()
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}

func TestLinkedMap_FirstOr(t *testing.T) {
	m := NewLinkedMap[int, int]()
	v := m.FirstOr(10)
	assert.Equal(t, 10, v)
}

func TestLinkedMap_Last(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(2, 2)
	m.Set(0, 0)
	m.Set(1, 1)
	v, ok := m.Last()
	assert.True(t, ok)
	assert.Equal(t, 1, v)
}

func TestLinkedMap_LastOr(t *testing.T) {
	m := NewLinkedMap[int, int]()
	v := m.LastOr(10)
	assert.Equal(t, 10, v)
}

func TestLinkedMap_Keys(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.Equal(t, []int{0, 1, 2}, m.Keys())
}

func TestLinkedMap_Values(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.Equal(t, []int{0, 1, 2}, m.Values())
}

func TestLinkedMap_Clear(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	m.Clear()
	assert.True(t, m.IsEmpty())
}

func TestLinkedMap_ContainsKey(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.ContainsKey(0))
}

func TestLinkedMap_Contains(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.Contains(0))
}

func TestLinkedMap_ContainsWhere(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	assert.True(t, m.ContainsWhere(func(value int) bool {
		return value == 2
	}))
}

func TestLinkedMap_Each(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	items := []int{}
	m.Each(func(key, value int) bool {
		items = append(items, value)
		return value < 1
	})
	assert.Equal(t, []int{0, 1}, items)
}

func TestLinkedMap_ToJSON(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	jsonBytes, err := m.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `{"entries":{"0":0,"1":1,"2":2},"keys":[0,1,2]}`, string(jsonBytes))
}

func TestLinkedMap_MarshalJSON(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	jsonBytes, err := json.Marshal(m)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"entries":{"0":0,"1":1,"2":2},"keys":[0,1,2]}`, string(jsonBytes))
}

func TestLinkedMap_UnmarshalJSON(t *testing.T) {
	m := NewLinkedMap[int, int]()
	err := json.Unmarshal([]byte(`{"entries":{"0":0,"1":1,"2":2},"keys":[2,0,1]}`), m)
	assert.Nil(t, err)
	assert.EqualValues(t, []int{2, 0, 1}, m.Keys())
	assert.EqualValues(t, map[int]int{
		0: 0, 1: 1, 2: 2,
	}, m.ToMap())
}

func TestLinkedMap_String(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(2, 2)
	m.Set(1, 1)
	str := m.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`LinkedMap\[int,\sint\]\(len=%d\)\{\n(\t\d+:\s\d+,\n)+\}`, m.Count()))
	assert.True(t, pattern.Match([]byte(str)))
}

func TestLinkedMap_Clone(t *testing.T) {
	m := NewLinkedMap[int, int]()
	m.Set(0, 0)
	m.Set(1, 1)
	m.Set(2, 2)
	m2 := m.Clone()
	assert.EqualValues(t, map[int]int{
		0: 0, 1: 1, 2: 2,
	}, m2.ToMap())
}
