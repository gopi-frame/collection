package list

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/gopi-frame/contract"
)

// NewList new list
func NewList[E any](values ...E) *List[E] {
	instance := new(List[E])
	instance.Push(values...)
	return instance
}

// List list
type List[E any] struct {
	sync.RWMutex
	items []E
}

// Count returns the size of the list
func (list *List[E]) Count() int64 {
	return int64(len(list.items))
}

// IsEmpty returns whether the list is empty.
func (list *List[E]) IsEmpty() bool {
	return list.Count() == 0
}

// IsNotEmpty returns whether the list is not empty.
func (list *List[E]) IsNotEmpty() bool {
	return !list.IsEmpty()
}

// Contains returns whether the list contains the specific element.
func (list *List[E]) Contains(value E) bool {
	return list.ContainsWhere(func(e E) bool {
		return reflect.DeepEqual(e, value)
	})
}

// ContainsWhere returns whether the list contains specific elements by callback.
func (list *List[E]) ContainsWhere(callback func(value E) bool) bool {
	return slices.ContainsFunc(list.items, callback)
}

// Push pushes elements into the list.
func (list *List[E]) Push(values ...E) {
	list.items = append(list.items, values...)
}

// Remove removes the specific element.
func (list *List[E]) Remove(value E) {
	list.RemoveWhere(func(item E) bool {
		return reflect.DeepEqual(value, item)
	})
}

// RemoveWhere removes specific elements by callback.
func (list *List[E]) RemoveWhere(callback func(item E) bool) {
	list.items = slices.DeleteFunc(list.items, callback)
}

// RemoveAt removes the element on the specific index.
func (list *List[E]) RemoveAt(index int) {
	list.items = slices.Delete(list.items, index, index+1)
}

// Clear clears the list.
func (list *List[E]) Clear() {
	list.items = []E{}
}

// Get returns the element on the specific index.
func (list *List[E]) Get(index int) E {
	return list.items[index]
}

// Set sets element on the specific index.
func (list *List[E]) Set(index int, value E) {
	list.items[index] = value
}

// First returns the first element of the list.
// it will return a zero value and false when the list is empty.
func (list *List[E]) First() (E, bool) {
	if len(list.items) == 0 {
		return *new(E), false
	}
	return list.items[0], true
}

// FirstOr returns the first element of the list, it will return the default value when the list is empty.
func (list *List[E]) FirstOr(value E) E {
	if v, ok := list.First(); ok {
		return v
	}
	return value
}

// FirstWhere returns the first element of the list which matches the callback.
// It will return a zero value and false when none matches the callback.
func (list *List[E]) FirstWhere(callback func(item E) bool) (E, bool) {
	for _, item := range list.items {
		if callback(item) {
			return item, true
		}
	}
	return *new(E), false
}

// FirstWhereOr returns the first element of the list which matches the callback.
// It will return the default value when none matches the callback.
func (list *List[E]) FirstWhereOr(callback func(item E) bool, value E) E {
	if v, found := list.FirstWhere(callback); found {
		return v
	}
	return value
}

// Last returns the last element of the list.
// It will return a zero value and false when the list is empty.
func (list *List[E]) Last() (E, bool) {
	length := len(list.items)
	if length == 0 {
		return *new(E), false
	}
	return list.items[length-1], true
}

// LastOr returns the last element of the list.
// It will return the default value when the list is empty.
func (list *List[E]) LastOr(value E) E {
	if v, ok := list.Last(); ok {
		return v
	}
	return value
}

// LastWhere returns the last element of the list which matches the callback.
// It will return a zero value and false when none matches the callback.
func (list *List[E]) LastWhere(callback func(item E) bool) (E, bool) {
	length := len(list.items)
	for index := range list.items {
		if value := list.items[length-index-1]; callback(value) {
			return value, true
		}
	}
	return *new(E), false
}

// LastWhereOr returns the last element of the list which matches the callback.
// It will return the default value when none matches the callback.
func (list *List[E]) LastWhereOr(callback func(item E) bool, value E) E {
	if v, ok := list.LastWhere(callback); ok {
		return v
	}
	return value
}

// Pop removes the last element of the list and returns it.
// It will return a zero value and false when the list is empty.
func (list *List[E]) Pop() (E, bool) {
	length := len(list.items)
	if length == 0 {
		return *new(E), false
	}
	value := list.items[length-1]
	list.items = list.items[:length-1]
	return value, true
}

// Shift removes the first element of the list and returns it.
// It will return a zero value and false when the list is empty.
func (list *List[E]) Shift() (E, bool) {
	if len(list.items) == 0 {
		return *new(E), false
	}
	value := list.items[0]
	list.items = list.items[1:]
	return value, true
}

// Unshift puts elements to the head of the list.
func (list *List[E]) Unshift(values ...E) {
	list.items = slices.Insert(list.items, 0, values...)
}

// IndexOf returns the index of the specific element.
func (list *List[E]) IndexOf(value E) int {
	return list.IndexOfWhere(func(item E) bool {
		return reflect.DeepEqual(value, item)
	})
}

// IndexOfWhere returns the index of the first element which matches the callback.
func (list *List[E]) IndexOfWhere(callback func(item E) bool) int {
	return slices.IndexFunc(list.items, callback)
}

// Sub returns the sub list with given range
func (list *List[E]) Sub(from, to int) *List[E] {
	return &List[E]{items: list.items[from:to]}
}

// Where returns the sub list with elements which matches the callback
func (list *List[E]) Where(callback func(item E) bool) *List[E] {
	l := &List[E]{}
	for _, item := range list.items {
		if callback(item) {
			l.items = append(l.items, item)
		}
	}
	return l
}

// Compact makes the list more compact
func (list *List[E]) Compact(callback func(a, b E) bool) {
	if callback == nil {
		callback = func(a, b E) bool {
			return reflect.DeepEqual(a, b)
		}
	}
	list.items = slices.CompactFunc(list.items, callback)
}

// Min returns the min element
func (list *List[E]) Min(callback func(a, b E) int) E {
	return slices.MinFunc(list.items, callback)
}

// Max returns the max element
func (list *List[E]) Max(callback func(a, b E) int) E {
	return slices.MaxFunc(list.items, callback)
}

// Sort sorts the list
func (list *List[E]) Sort(callback func(a, b E) int) {
	slices.SortFunc(list.items, callback)
}

// Chunk splits list into multiply parts by given size
func (list *List[E]) Chunk(size int) *List[*List[any]] {
	chunks := NewList[*List[any]]()
	chunk := NewList[any]()
	for _, item := range list.items {
		if len(chunk.items) < size {
			chunk.Push(item)
		} else {
			chunks.Push(chunk)
			chunk = NewList[any](item)
		}
	}
	chunks.Push(chunk)
	return chunks
}

// Each travers the list, if the callback returns false then break
func (list *List[E]) Each(callback func(index int, value E) bool) {
	for index, value := range list.items {
		if !callback(index, value) {
			break
		}
	}
}

// Reverse reverses the list
func (list *List[E]) Reverse() {
	slices.Reverse(list.items)
}

// Clone clones the list
func (list *List[E]) Clone() *List[E] {
	list.items = slices.Clone(list.items)
	return list
}

// String convert to string
func (list *List[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("List[%T](len=%d)", *new(E), list.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	for index, value := range list.items {
		str.WriteByte('\t')
		if v, ok := any(value).(contract.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		if index >= 4 {
			break
		}
	}
	if list.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}

// ToJSON converts to json
func (list *List[E]) ToJSON() ([]byte, error) {
	return json.Marshal(list.items)
}

// ToArray converts to array
func (list *List[E]) ToArray() []E {
	return list.items
}

// MarshalJSON implements [json.Marshaller]
func (list *List[E]) MarshalJSON() ([]byte, error) {
	return list.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (list *List[E]) UnmarshalJSON(data []byte) error {
	var items []E
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	list.items = items
	return nil
}
