package tree

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gopi-frame/contract"
)

// NewAVLTree new avl tree
func NewAVLTree[E any](comparator contract.Comparator[E], values ...E) *AVLTree[E] {
	tree := new(AVLTree[E])
	tree.comparator = comparator
	tree.Push(values...)
	return tree
}

// AVLTree avl tree
type AVLTree[E any] struct {
	sync.Mutex
	root       *avlNode[E]
	comparator contract.Comparator[E]
}

// Count returns the size of tree
func (t *AVLTree[E]) Count() int64 {
	return int64(len(t.root.inOrderRange()))
}

// IsEmpty returns whether the tree is empty
func (t *AVLTree[E]) IsEmpty() bool {
	return t.Count() == 0
}

// IsNotEmpty returns whether the tree is not empty
func (t *AVLTree[E]) IsNotEmpty() bool {
	return t.Count() > 0
}

// Contains returns whether the tree contains the specific element
func (t *AVLTree[E]) Contains(value E) bool {
	if t.root == nil {
		return false
	}
	if t.root.find(value, t.comparator) == nil {
		return false
	}
	return true
}

// Push pushes elements into the tree
func (t *AVLTree[E]) Push(values ...E) {
	for _, value := range values {
		t.root = t.root.insert(value, t.comparator)
	}
}

// Remove removes the specific element from the tree
func (t *AVLTree[E]) Remove(value E) {
	if t.root == nil {
		return
	}
	t.root = t.root.remove(value, t.comparator)
}

// Clear clears the tree
func (t *AVLTree[E]) Clear() {
	t.root = nil
}

// First returns the first element of the tree.
// It returns zero value and false when the tree is empty.
func (t *AVLTree[E]) First() (E, bool) {
	if t.root == nil {
		return *new(E), false
	}
	return t.root.min().value, true
}

// FirstOr returns the first element of the tree or the default value if the tree is empty
func (t *AVLTree[E]) FirstOr(value E) E {
	if t.root == nil {
		return value
	}
	return t.root.min().value
}

// Last returns the last element of the tree.
// It returns zero value and false when the tree is empty
func (t *AVLTree[E]) Last() (E, bool) {
	if t.root == nil {
		return *new(E), false
	}
	return t.root.max().value, true
}

// LastOr returns the last element of the tree or the default value if the tree is empty
func (t *AVLTree[E]) LastOr(value E) E {
	if t.root == nil {
		return value
	}
	return t.root.max().value
}

// Each runs callback for each element, it breaks when callback returns false
func (t *AVLTree[E]) Each(callback func(value E) bool) {
	for _, node := range t.root.inOrderRange() {
		if !callback(node.value) {
			break
		}
	}
}

// Clone clones the tree
func (t *AVLTree[E]) Clone() *AVLTree[E] {
	tt := NewAVLTree(t.comparator, t.ToArray()...)
	return tt
}

// ToArray converts to array
func (t *AVLTree[E]) ToArray() []E {
	nodes := t.root.inOrderRange()
	values := make([]E, 0, len(nodes))
	for _, node := range nodes {
		values = append(values, node.value)
	}
	return values
}

// ToJSON converts to json
func (t *AVLTree[E]) ToJSON() ([]byte, error) {
	return json.Marshal(t.ToArray())
}

// MarshalJSON implements [json.Marshaller]
func (t *AVLTree[E]) MarshalJSON() ([]byte, error) {
	return t.ToJSON()
}

// UnmarshalJSON implements [json.UnmarshalJSON]
func (t *AVLTree[E]) UnmarshalJSON(data []byte) error {
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	t.Clear()
	t.Push(values...)
	return nil
}

// String converts to string
func (t *AVLTree[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("AVLTree[%T](len=%d)", *new(E), t.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	items := t.ToArray()
	for index, item := range items {
		str.WriteByte('\t')
		if v, ok := any(item).(contract.Stringable); ok {
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
	if len(items) > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
