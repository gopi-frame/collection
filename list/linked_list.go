package list

import (
	listlib "container/list"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/gopi-frame/contract"
	"github.com/gopi-frame/exception"
)

// NewLinkedList new linked list
func NewLinkedList[E any](values ...E) *LinkedList[E] {
	instance := new(LinkedList[E])
	instance.Push(values...)
	return instance
}

// LinkedList linked list
type LinkedList[E any] struct {
	sync.RWMutex
	list *listlib.List
}

func (l *LinkedList[E]) init() {
	if l.list == nil {
		l.list = listlib.New()
	}
}

// Count returns the size of the list
func (l *LinkedList[E]) Count() int64 {
	l.init()
	return int64(l.list.Len())
}

// IsEmpty returns whether the list is empty.
func (l *LinkedList[E]) IsEmpty() bool {
	l.init()
	return l.Count() == 0
}

// IsNotEmpty returns whether the list is not empty.
func (l *LinkedList[E]) IsNotEmpty() bool {
	l.init()
	return !l.IsEmpty()
}

// Contains returns whether the list contains the specific element.
func (l *LinkedList[E]) Contains(value E) bool {
	l.init()
	return l.ContainsWhere(func(item E) bool {
		return reflect.DeepEqual(item, value)
	})
}

// ContainsWhere returns whether the list contains specific elements by callback.
func (l *LinkedList[E]) ContainsWhere(callback func(value E) bool) bool {
	l.init()
	for e := l.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			return true
		}
	}
	return false
}

// Push pushes elements into the list.
func (l *LinkedList[E]) Push(values ...E) {
	l.init()
	for _, value := range values {
		l.list.PushBack(value)
	}
}

// Remove removes the specific element.
func (l *LinkedList[E]) Remove(value E) {
	l.RemoveWhere(func(item E) bool {
		return reflect.DeepEqual(item, value)
	})
}

// RemoveWhere removes specific elements by callback.
func (l *LinkedList[E]) RemoveWhere(callback func(item E) bool) {
	l.init()
	var next *listlib.Element
	for e := l.list.Front(); e != nil; e = next {
		next = e.Next()
		if callback(e.Value.(E)) {
			l.list.Remove(e)
		}
	}
}

// RemoveAt removes the element on the specific index.
func (l *LinkedList[E]) RemoveAt(index int) {
	l.init()
	var next *listlib.Element
	for e, i := l.list.Front(), 0; e != nil; e, i = next, i+1 {
		next = e.Next()
		if i == index {
			l.list.Remove(e)
			break
		}
	}
}

// Clear clears the list.
func (l *LinkedList[E]) Clear() {
	l.init()
	l.list.Init()
}

// Get returns the element on the specific index.
func (l *LinkedList[E]) Get(index int) E {
	l.init()
	if index < 0 || index >= l.list.Len() {
		panic(exception.NewRangeException(0, l.list.Len()-1))
	}
	for i, e := 0, l.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i == index {
			return e.Value.(E)
		}
	}
	return *new(E)
}

// Set sets element on the specific index.
func (l *LinkedList[E]) Set(index int, value E) {
	l.init()
	for i, e := 0, l.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i == index {
			e.Value = value
		}
	}
}

// First returns the first element of the list.
// it will return a zero value and false when the list is empty.
func (l *LinkedList[E]) First() (E, bool) {
	l.init()
	if l.list.Len() == 0 {
		return *new(E), false
	}
	return l.list.Front().Value.(E), true
}

// FirstOr returns the first element of the list, it will return the default value when the list is empty.
func (l *LinkedList[E]) FirstOr(value E) E {
	l.init()
	if l.list.Len() == 0 {
		return value
	}
	return l.list.Front().Value.(E)
}

// FirstWhere returns the first element of the list which matches the callback.
// It will return a zero value and false when none matches the callback.
func (l *LinkedList[E]) FirstWhere(callback func(item E) bool) (E, bool) {
	l.init()
	for e := l.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			return e.Value.(E), true
		}
	}
	return *new(E), false
}

// FirstWhereOr returns the first element of the list which matches the callback.
// It will return the default value when none matches the callback.
func (l *LinkedList[E]) FirstWhereOr(callback func(item E) bool, value E) E {
	l.init()
	for e := l.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			return e.Value.(E)
		}
	}
	return value
}

// Last returns the last element of the list.
// It will return a zero value and false when the list is empty.
func (l *LinkedList[E]) Last() (E, bool) {
	l.init()
	if l.list.Len() == 0 {
		return *new(E), false
	}
	return l.list.Back().Value.(E), true
}

// LastOr returns the last element of the list.
// It will return the default value when the list is empty.
func (l *LinkedList[E]) LastOr(value E) E {
	l.init()
	if l.list.Back() == nil {
		return value
	}
	return l.list.Back().Value.(E)
}

// LastWhere returns the last element of the list which matches the callback.
// It will return a zero value and false when none matches the callback.
func (l *LinkedList[E]) LastWhere(callback func(item E) bool) (E, bool) {
	l.init()
	for e := l.list.Back(); e != nil; e = e.Prev() {
		if callback(e.Value.(E)) {
			return e.Value.(E), true
		}
	}
	return *new(E), false
}

// LastWhereOr returns the last element of the list which matches the callback.
// It will return the default value when none matches the callback.
func (l *LinkedList[E]) LastWhereOr(callback func(item E) bool, value E) E {
	l.init()
	if v, ok := l.LastWhere(callback); ok {
		return v
	}
	return value
}

// Pop removes the last element of the list and returns it.
// It will return a zero value and false when the list is empty.
func (l *LinkedList[E]) Pop() (E, bool) {
	l.init()
	if l.list.Len() == 0 {
		return *new(E), false
	}
	item := l.list.Back()
	l.list.Remove(item)
	return item.Value.(E), true
}

// Shift removes the first element of the list and returns it.
// It will return a zero value and false when the list is empty.
func (l *LinkedList[E]) Shift() (E, bool) {
	l.init()
	if l.list.Len() == 0 {
		return *new(E), false
	}
	item := l.list.Front()
	l.list.Remove(item)
	return item.Value.(E), true
}

// Unshift puts elements to the head of the list.
func (l *LinkedList[E]) Unshift(values ...E) {
	l.init()
	for _, value := range values {
		l.list.PushFront(value)
	}
}

// IndexOf returns the index of the specific element.
func (l *LinkedList[E]) IndexOf(value E) int {
	l.init()
	return l.IndexOfWhere(func(item E) bool {
		return reflect.DeepEqual(item, value)
	})
}

// IndexOfWhere returns the index of the first element which matches the callback.
func (l *LinkedList[E]) IndexOfWhere(callback func(item E) bool) int {
	l.init()
	for i, e := 0, l.list.Front(); e != nil; i, e = i+1, e.Next() {
		if callback(e.Value.(E)) {
			return i
		}
	}
	return -1
}

// Sub returns the sub list with given range
func (l *LinkedList[E]) Sub(from, to int) *LinkedList[E] {
	l.init()
	linked := NewLinkedList[E]()
	for i, e := 0, l.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i < from {
			continue
		} else if i >= from && i < to {
			linked.Push(e.Value.(E))
		} else {
			break
		}
	}
	return linked
}

// Where returns the sub list with elements which matches the callback
func (l *LinkedList[E]) Where(callback func(item E) bool) *LinkedList[E] {
	l.init()
	linked := &LinkedList[E]{}
	for e := l.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			linked.Push(e.Value.(E))
		}
	}
	return linked
}

// Compact makes the list more compact
func (l *LinkedList[E]) Compact(callback func(a, b E) bool) {
	l.init()
	if l.list.Len() < 2 {
		return
	}
	if callback == nil {
		callback = func(a, b E) bool {
			return reflect.DeepEqual(a, b)
		}
	}
	var next *listlib.Element
	for e := l.list.Front().Next(); e != nil; e = next {
		next = e.Next()
		if callback(e.Value.(E), e.Prev().Value.(E)) {
			l.list.Remove(e)
		}
	}
}

// Min returns the min element
func (l *LinkedList[E]) Min(callback func(a, b E) int) E {
	l.init()
	return slices.MinFunc(l.ToArray(), callback)
}

// Max returns the max element
func (l *LinkedList[E]) Max(callback func(a, b E) int) E {
	l.init()
	return slices.MaxFunc(l.ToArray(), callback)
}

// Sort sorts the list
func (l *LinkedList[E]) Sort(callback func(a, b E) int) {
	l.init()
	var newList = listlib.New()
	for e := l.list.Front(); e != nil; e = e.Next() {
		node := newList.Front()
		for node != nil {
			if callback(e.Value.(E), node.Value.(E)) < 0 {
				newList.InsertBefore(e.Value, node)
				break
			}
			node = node.Next()
		}
		if node == nil {
			newList.PushBack(e.Value)
		}
	}
	l.list = newList
}

// Chunk splits list into multiply parts by given size
func (l *LinkedList[E]) Chunk(size int) *LinkedList[*LinkedList[any]] {
	l.init()
	chunks := NewLinkedList[*LinkedList[any]]()
	chunk := NewLinkedList[any]()
	for e := l.list.Front(); e != nil; e = e.Next() {
		if chunk.list.Len() < size {
			chunk.Push(e.Value.(E))
		} else {
			chunks.Push(chunk)
			chunk = NewLinkedList(e.Value)
		}
	}
	chunks.Push(chunk)
	return chunks
}

// Each travers the list, if the callback returns false then break
func (l *LinkedList[E]) Each(callback func(index int, value E) bool) {
	l.init()
	for e, i := l.list.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		if !callback(i, e.Value.(E)) {
			break
		}
	}
}

// Reverse reverses the list
func (l *LinkedList[E]) Reverse() {
	l.init()
	var next *listlib.Element
	for e := l.list.Front(); e != nil; e = next {
		next = e.Next()
		l.list.PushFront(e.Value)
		l.list.Remove(e)
	}
}

// Clone clones the list
func (l *LinkedList[E]) Clone() *LinkedList[E] {
	l.init()
	linked := &LinkedList[E]{}
	for e := l.list.Front(); e != nil; e = e.Next() {
		linked.Push(e.Value.(E))
	}
	return linked
}

// String convert to string
func (l *LinkedList[E]) String() string {
	l.init()
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedList[%T](len=%d)", *new(E), l.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	l.Each(func(index int, value E) bool {
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
	if l.list.Len() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}

// ToJSON converts to json
func (l *LinkedList[E]) ToJSON() ([]byte, error) {
	l.init()
	return json.Marshal(l.ToArray())
}

// ToArray converts to array
func (l *LinkedList[E]) ToArray() []E {
	l.init()
	var items []E
	for e := l.list.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(E))
	}
	return items
}

// MarshalJSON implements [json.Marshaller]
func (l *LinkedList[E]) MarshalJSON() ([]byte, error) {
	l.init()
	return l.ToJSON()
}

// UnmarshalJSON implements [json.Unmarshaller]
func (l *LinkedList[E]) UnmarshalJSON(data []byte) error {
	l.init()
	items := []E{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		l.list.PushBack(item)
	}
	return nil
}
