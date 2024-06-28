package set

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/contract"
)

// NewLinkedSet creates a new linked hash set
func NewLinkedSet[E comparable](values ...E) *LinkedSet[E] {
	set := new(LinkedSet[E])
	set.elements = map[E]struct{}{}
	set.link = list.NewLinkedList[E]()
	set.Push(values...)
	return set
}

// LinkedSet linked hash set
type LinkedSet[E comparable] struct {
	sync.RWMutex
	elements map[E]struct{}
	link     *list.LinkedList[E]
}

func (s *LinkedSet[E]) Count() int64 {
	return s.link.Count()
}

func (s *LinkedSet[E]) IsEmpty() bool {
	return s.Count() == 0
}

func (s *LinkedSet[E]) IsNotEmpty() bool {
	return !s.IsEmpty()
}

func (s *LinkedSet[E]) Contains(value E) bool {
	_, contains := s.elements[value]
	return contains
}

func (s *LinkedSet[E]) ContainsWhere(callback func(E) bool) bool {
	return s.link.ContainsWhere(callback)
}

func (s *LinkedSet[E]) Push(values ...E) {
	for _, value := range values {
		if s.Contains(value) {
			continue
		}
		s.elements[value] = struct{}{}
		s.link.Push(value)
	}
}

func (s *LinkedSet[E]) Remove(value E) {
	s.RemoveWhere(func(e E) bool {
		return e == value
	})
}

func (s *LinkedSet[E]) RemoveWhere(callback func(E) bool) {
	s.link = s.link.Where(func(item E) bool {
		return !callback(item)
	})
	s.elements = make(map[E]struct{})
	s.link.Each(func(index int, value E) bool {
		s.elements[value] = struct{}{}
		return true
	})
}

func (s *LinkedSet[E]) Clear() {
	s.link.Clear()
}

func (s *LinkedSet[E]) Each(callback func(int, E) bool) {
	s.link.Each(callback)
}

func (s *LinkedSet[E]) Clone() *LinkedSet[E] {
	return NewLinkedSet(s.ToArray()...)
}

func (s *LinkedSet[E]) ToArray() []E {
	return s.link.ToArray()
}

func (s *LinkedSet[E]) ToJSON() ([]byte, error) {
	return json.Marshal(s.ToArray())
}

func (s *LinkedSet[E]) MarshalJSON() ([]byte, error) {
	return s.ToJSON()
}

func (s *LinkedSet[E]) UnmarshalJSON(data []byte) error {
	items := list.NewLinkedList[E]()
	err := json.Unmarshal(data, items)
	if err != nil {
		return err
	}
	s.link = items
	return nil
}

func (s *LinkedSet[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedSet[%T](len=%d)", *new(E), s.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	s.link.Each(func(index int, value E) bool {
		str.WriteByte('\t')
		if v, ok := any(value).(contract.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		return index < 4
	})
	if s.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
