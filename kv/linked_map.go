package kv

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/contract"
)

type jsonObject[K comparable, V any] struct {
	Entries map[K]V `json:"entries"`
	Keys    []K     `json:"keys"`
}

// NewLinkedMap new linked map
func NewLinkedMap[K comparable, V any]() *LinkedMap[K, V] {
	m := new(LinkedMap[K, V])
	m.Map = NewMap[K, V]()
	m.keys = list.NewLinkedList[K]()
	return m
}

// LinkedMap linked map
type LinkedMap[K comparable, V any] struct {
	sync.RWMutex
	*Map[K, V]
	keys *list.LinkedList[K]
}

func (m *LinkedMap[K, V]) Set(key K, value V) {
	m.Map.Set(key, value)
	m.keys.Push(key)
}

func (m *LinkedMap[K, V]) Remove(key K) {
	m.Map.Remove(key)
	m.keys.Remove(key)
}

func (m *LinkedMap[K, V]) First() (V, bool) {
	if len(m.items) == 0 {
		return *new(V), false
	}
	k, _ := m.keys.First()
	v, ok := m.items[k]
	return v, ok
}

func (m *LinkedMap[K, V]) FirstOr(value V) V {
	if v, ok := m.First(); ok {
		return v
	}
	return value
}

func (m *LinkedMap[K, V]) Last() (V, bool) {
	if len(m.items) == 0 {
		return *new(V), false
	}
	k, _ := m.keys.Last()
	v, ok := m.items[k]
	return v, ok
}

func (m *LinkedMap[K, V]) LastOr(value V) V {
	if v, ok := m.Last(); ok {
		return v
	}
	return value
}

func (m *LinkedMap[K, V]) Keys() []K {
	keys := []K{}
	m.keys.Each(func(index int, value K) bool {
		keys = append(keys, value)
		return true
	})
	return keys
}

func (m *LinkedMap[K, V]) Values() []V {
	values := []V{}
	m.keys.Each(func(index int, value K) bool {
		values = append(values, m.items[value])
		return true
	})
	return values
}

func (m *LinkedMap[K, V]) Clear() {
	m.items = make(map[K]V)
	m.keys.Clear()
}

func (m *LinkedMap[K, V]) ContainsKey(key K) bool {
	for k := range m.items {
		if k == key {
			return true
		}
	}
	return false
}

func (m *LinkedMap[K, V]) Reverse() *LinkedMap[K, V] {
	m.keys.Reverse()
	return m
}

func (m *LinkedMap[K, V]) Each(callback func(key K, value V) bool) {
	m.keys.Each(func(index int, value K) bool {
		return callback(value, m.items[value])
	})
}

func (m *LinkedMap[K, V]) ToJSON() ([]byte, error) {
	return json.Marshal(jsonObject[K, V]{
		Entries: m.ToMap(),
		Keys:    m.keys.ToArray(),
	})
}

func (m *LinkedMap[K, V]) MarshalJSON() ([]byte, error) {
	return m.ToJSON()
}

func (m *LinkedMap[K, V]) UnmarshalJSON(data []byte) error {
	var container = new(jsonObject[K, V])
	err := json.Unmarshal(data, container)
	if err != nil {
		return err
	}
	m.Map = NewMap[K, V]()
	m.keys = list.NewLinkedList(container.Keys...)
	m.keys.Each(func(index int, value K) bool {
		m.Map.Set(value, container.Entries[value])
		return true
	})
	return nil
}

func (m *LinkedMap[K, V]) ToMap() map[K]V {
	return m.items
}

func (m *LinkedMap[K, V]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedMap[%T, %T](len=%d)", *new(K), *new(V), m.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	keys := m.Keys()
	for _, key := range keys {
		str.WriteByte('\t')
		if k, ok := any(key).(contract.Stringable); ok {
			str.WriteString(k.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", key))
		}
		str.WriteByte(':')
		str.WriteByte(' ')
		value, _ := m.Map.Get(key)
		if v, ok := any(value).(contract.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
	}
	str.WriteByte('}')
	return str.String()
}

func (m *LinkedMap[K, V]) Clone() *LinkedMap[K, V] {
	mm := NewLinkedMap[K, V]()
	m.keys.Each(func(index int, key K) bool {
		mm.Set(key, m.items[key])
		return true
	})
	return m
}
