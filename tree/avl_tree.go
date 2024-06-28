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

// AVLTree avl true
type AVLTree[E any] struct {
	sync.Mutex
	root       *avlNode[E]
	comparator contract.Comparator[E]
}

func (t *AVLTree[E]) Count() int64 {
	return int64(len(t.root.inOrderRange()))
}

func (t *AVLTree[E]) IsEmpty() bool {
	return t.Count() == 0
}

func (t *AVLTree[E]) IsNotEmpty() bool {
	return t.Count() > 0
}

func (t *AVLTree[E]) Contains(value E) bool {
	if t.root == nil {
		return false
	}
	if t.root.find(value, t.comparator) == nil {
		return false
	}
	return true
}

func (t *AVLTree[E]) Push(values ...E) {
	for _, value := range values {
		t.root = t.root.insert(value, t.comparator)
	}
}

func (t *AVLTree[E]) Remove(value E) {
	if t.root == nil {
		return
	}
	t.root = t.root.remove(value, t.comparator)
}

func (t *AVLTree[E]) Clear() {
	t.root = nil
}

func (t *AVLTree[E]) First() (E, bool) {
	if t.root == nil {
		return *new(E), false
	}
	return t.root.min().value, true
}

func (t *AVLTree[E]) FirstOr(value E) E {
	if t.root == nil {
		return value
	}
	return t.root.min().value
}

func (t *AVLTree[E]) Last() (E, bool) {
	if t.root == nil {
		return *new(E), false
	}
	return t.root.max().value, true
}

func (t *AVLTree[E]) LastOr(value E) E {
	if t.root == nil {
		return value
	}
	return t.root.max().value
}

func (t *AVLTree[E]) Each(callback func(value E) bool) {
	for _, node := range t.root.inOrderRange() {
		if !callback(node.value) {
			break
		}
	}
}

func (t *AVLTree[E]) Clone() *AVLTree[E] {
	avltree := NewAVLTree(t.comparator, t.ToArray()...)
	return avltree
}

func (t *AVLTree[E]) ToArray() []E {
	nodes := t.root.inOrderRange()
	values := make([]E, 0, len(nodes))
	for _, node := range nodes {
		values = append(values, node.value)
	}
	return values
}

func (t *AVLTree[E]) ToJSON() ([]byte, error) {
	return json.Marshal(t.ToArray())
}

func (t *AVLTree[E]) MarshalJSON() ([]byte, error) {
	return t.ToJSON()
}

func (t *AVLTree[E]) UnmarshalJSON(data []byte) error {
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	t.Clear()
	t.Push(values...)
	return nil
}

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
