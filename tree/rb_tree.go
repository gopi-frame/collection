package tree

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gopi-frame/contract"
)

// NewRBTree new rb tree
func NewRBTree[E any](comparator contract.Comparator[E], values ...E) *RBTree[E] {
	tree := new(RBTree[E])
	tree.comparator = comparator
	tree.Push(values...)
	return tree
}

// RBTree red black tree
type RBTree[E any] struct {
	sync.Mutex
	root       *rbNode[E]
	comparator contract.Comparator[E]
}

func (t *RBTree[E]) Count() int64 {
	return int64(len(t.root.inOrderRange()))
}

func (t *RBTree[E]) IsEmpty() bool {
	return t.Count() == 0
}

func (t *RBTree[E]) IsNotEmpty() bool {
	return t.Count() > 0
}

func (t *RBTree[E]) Contains(value E) bool {
	if t.root == nil {
		return false
	}
	if t.root.find(value, t.comparator) == nil {
		return false
	}
	return true
}

func (t *RBTree[E]) Push(values ...E) *RBTree[E] {
	for _, value := range values {
		t.root = t.root.insert(value, t.comparator)
		t.root.color = black
	}
	return t
}

func (t *RBTree[E]) Remove(value E) *RBTree[E] {
	if t.root == nil {
		return t
	}
	if t.root.find(value, t.comparator) == nil {
		return t
	}
	if t.root.left.isBlack() && t.root.right.isBlack() {
		t.root.color = red
	}
	t.root = t.root.remove(value, t.comparator)
	if t.root.isRed() {
		t.root.color = black
	}
	return t
}

func (t *RBTree[E]) Clear() *RBTree[E] {
	t.root = nil
	return t
}

func (t *RBTree[E]) Comparator() contract.Comparator[E] {
	return t.comparator
}

func (t *RBTree[E]) First() (E, bool) {
	if t.root == nil {
		return *new(E), false
	}
	value := t.root.min().value
	return value, true
}

func (t *RBTree[E]) FirstOr(value E) E {
	if t.root == nil {
		return value
	}
	v := t.root.min().value
	return v
}

func (t *RBTree[E]) Last() (E, bool) {
	if t.root == nil {
		return *new(E), false
	}
	value := t.root.max().value
	return value, true
}

func (t *RBTree[E]) LastOr(value E) E {
	if t.root == nil {
		return value
	}
	v := t.root.max().value
	return v
}

func (t *RBTree[E]) Each(callback func(value E) bool) {
	for _, node := range t.root.inOrderRange() {
		if !callback(node.value) {
			break
		}
	}
}

func (t *RBTree[E]) Clone() *RBTree[E] {
	rbTree := NewRBTree(t.comparator, t.ToArray()...)
	return rbTree
}

func (t *RBTree[E]) ToArray() []E {
	nodes := t.root.inOrderRange()
	values := make([]E, 0, len(nodes))
	for _, node := range nodes {
		values = append(values, node.value)
	}
	return values
}

func (t *RBTree[E]) ToJSON() ([]byte, error) {
	return json.Marshal(t.ToArray())
}

func (t *RBTree[E]) MarshalJSON() ([]byte, error) {
	return t.ToJSON()
}

func (t *RBTree[E]) UnmarshalJSON(data []byte) error {
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	t.Clear().Push(values...)
	return nil
}

func (t *RBTree[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("RBTree[%T](len=%d)", *new(E), t.Count()))
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
