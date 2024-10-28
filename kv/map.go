package kv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gopi-frame/contract"
)

// NewMap new map
func NewMap[K comparable, V any]() *Map[K, V] {
	m := new(Map[K, V])
	m.items = make(map[K]V)
	return m
}

// NewFromMap new from map
func NewFromMap[K comparable, V any](m map[K]V) *Map[K, V] {
	mm := NewMap[K, V]()
	mm.items = m
	return mm
}

// Map map
type Map[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

// Count returns the size of map
func (m *Map[K, V]) Count() int64 {
	return int64(len(m.items))
}

// IsEmpty returns whether the map is empty
func (m *Map[K, V]) IsEmpty() bool {
	return m.Count() == 0
}

// IsNotEmpty returns whether the map is not empty
func (m *Map[K, V]) IsNotEmpty() bool {
	return !m.IsEmpty()
}

// Get gets element by specific key.
// A zero value and false will be returned when the given key is not exist
func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.items[key]
	return v, ok
}

// GetOr gets element by specific key
// The default will be return
func (m *Map[K, V]) GetOr(key K, value V) V {
	v, ok := m.items[key]
	if ok {
		return v
	}
	return value
}

// Set sets element to the specific key
func (m *Map[K, V]) Set(key K, value V) {
	m.items[key] = value
}

// Remove removes the element of specific key
func (m *Map[K, V]) Remove(key K) {
	delete(m.items, key)
}

// Keys returns all keys
func (m *Map[K, V]) Keys() []K {
	var keys []K
	for key := range m.items {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values
func (m *Map[K, V]) Values() []V {
	var values []V
	for _, value := range m.items {
		values = append(values, value)
	}
	return values
}

// Clear clears the map
func (m *Map[K, V]) Clear() {
	m.items = make(map[K]V)
}

// ContainsKey returns whether the map contains the specific key
func (m *Map[K, V]) ContainsKey(key K) bool {
	for k := range m.items {
		if k == key {
			return true
		}
	}
	return false
}

// Contains returns whether the map contains the specific value
func (m *Map[K, V]) Contains(value V) bool {
	return m.ContainsWhere(func(v V) bool {
		return reflect.DeepEqual(v, value)
	})
}

// ContainsWhere returns whether the map contains specific values through callback
func (m *Map[K, V]) ContainsWhere(callback func(value V) bool) bool {
	for _, v := range m.items {
		if callback(v) {
			return true
		}
	}
	return false
}

// Each ranges the map by callback, it will break the loop when the callback returns false
func (m *Map[K, V]) Each(callback func(key K, value V) bool) {
	for key, value := range m.items {
		if !callback(key, value) {
			break
		}
	}
}

// ToJSON converts the map to json bytes
func (m *Map[K, V]) ToJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// MarshalJSON implements [json.Marshaller]
func (m *Map[K, V]) MarshalJSON() ([]byte, error) {
	return m.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	values := map[K]V{}
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	m.items = values
	return nil
}

// ToMap converts to map
func (m *Map[K, V]) ToMap() map[K]V {
	return m.items
}

func (m *Map[K, V]) FromMap(items map[K]V) {
	m.items = items
}

// String converts to string
func (m *Map[K, V]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("Map[%T, %T](len=%d)", *new(K), *new(V), m.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	for k, v := range m.items {
		str.WriteByte('\t')
		if key, ok := any(k).(contract.Stringable); ok {
			str.WriteString(key.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", k))
		}
		str.WriteByte(':')
		str.WriteByte(' ')
		if value, ok := any(v).(contract.Stringable); ok {
			str.WriteString(value.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", v))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
	}
	str.WriteByte('}')
	return str.String()
}

// Clone clone a new map
func (m *Map[K, V]) Clone() *Map[K, V] {
	newMap := NewMap[K, V]()
	for key, value := range m.items {
		newMap.Set(key, value)
	}
	return newMap
}
