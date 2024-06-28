package set

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// NewSet new set
func NewSet[E comparable](values ...E) *Set[E] {
	set := &Set[E]{
		elements: make(map[E]struct{}),
	}
	for _, value := range values {
		set.elements[value] = struct{}{}
	}
	return set
}

// Set hash set
type Set[E comparable] struct {
	sync.RWMutex
	elements map[E]struct{}
}

// Count count
func (s *Set[E]) Count() int64 {
	return int64(len(s.elements))
}

// IsEmpty is empty
func (s *Set[E]) IsEmpty() bool {
	return s.Count() == 0
}

// IsNotEmpty is not empty
func (s *Set[E]) IsNotEmpty() bool {
	return !s.IsEmpty()
}

// Contains contains
func (s *Set[E]) Contains(value E) bool {
	_, contains := s.elements[value]
	return contains
}

// ContainsWhere comtains where
func (s *Set[E]) ContainsWhere(callback func(E) bool) bool {
	for item := range s.elements {
		if callback(item) {
			return true
		}
	}
	return false
}

// Push push
func (s *Set[E]) Push(values ...E) {
	for _, value := range values {
		if s.Contains(value) {
			continue
		}
		s.elements[value] = struct{}{}
	}
}

// Remove remove
func (s *Set[E]) Remove(value E) {
	delete(s.elements, value)
}

// RemoveWhere remove where
func (s *Set[E]) RemoveWhere(callback func(E) bool) {
	items := map[E]struct{}{}
	for item := range s.elements {
		if callback(item) {
			continue
		}
		items[item] = struct{}{}
	}
	s.elements = items
}

// Each each
func (s *Set[E]) Each(callback func(_ int, item E) bool) {
	for item := range s.elements {
		if !callback(-1, item) {
			break
		}
	}
}

// Clear clear
func (s *Set[E]) Clear() {
	s.elements = map[E]struct{}{}
}

// Clone clone
func (s *Set[E]) Clone() *Set[E] {
	return &Set[E]{
		elements: s.elements,
	}
}

// ToArray to array
func (s *Set[E]) ToArray() []E {
	values := []E{}
	for item := range s.elements {
		values = append(values, item)
	}
	return values
}

// ToJSON to json
func (s *Set[E]) ToJSON() ([]byte, error) {
	return json.Marshal(s.ToArray())
}

func (s *Set[E]) MarshalJSON() ([]byte, error) {
	return s.ToJSON()
}

func (s *Set[E]) UnmarshalJSON(data []byte) error {
	var items = []E{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	s.Clear()
	s.Push(items...)
	return nil
}

func (s *Set[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("Set[%T](len=%d)", *new(E), len(s.elements)))
	str.WriteByte('{')
	str.WriteByte('\n')
	index := 0
	for item := range s.elements {
		index++
		str.WriteByte('\t')
		if v, ok := any(item).(fmt.Stringer); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", item))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		if index >= 4 {
			break
		}
	}
	if len(s.elements) > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
