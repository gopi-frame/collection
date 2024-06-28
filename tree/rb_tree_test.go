package tree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBTree_Count(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3)
	assert.Equal(t, int64(3), tree.Count())
}

func TestRBTree_IsEmpty(t *testing.T) {
	tree := NewRBTree(_cmp{})
	assert.True(t, tree.IsEmpty())
}

func TestRBTree_IsNotEmpty(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 4)
	assert.True(t, tree.IsNotEmpty())
}

func TestRBTree_Contains(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3)
	ok := tree.Contains(1)
	assert.True(t, ok)
}

func TestRBTree_Remove(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3)
	ok := tree.Contains(1)
	assert.True(t, ok)
	tree.Remove(1)
	ok = tree.Contains(1)
	assert.False(t, ok)
}

func TestRBTree_Clear(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3)
	assert.False(t, tree.IsEmpty())
	tree.Clear()
	assert.True(t, tree.IsEmpty())
}

func TestRBTree_First(t *testing.T) {
	tree := NewRBTree(_cmp{}, 2, 1, 3, 0)
	v, ok := tree.First()
	assert.True(t, ok)
	assert.Equal(t, 0, v)
}

func TestRBTree_FirstOr(t *testing.T) {
	tree := NewRBTree(_cmp{})
	v := tree.FirstOr(1)
	assert.Equal(t, 1, v)
}

func TestRBTree_Last(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 4, 2)
	v, ok := tree.Last()
	assert.True(t, ok)
	assert.Equal(t, 4, v)
}

func TestRBTree_LastOr(t *testing.T) {
	tree := NewRBTree(_cmp{})
	v := tree.LastOr(1)
	assert.Equal(t, 1, v)
}

func TestRBTree_Each(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	items := []int{}
	tree.Each(func(value int) bool {
		items = append(items, value)
		return value < 2
	})
	assert.Equal(t, []int{1, 2}, items)
}

func TestRBTree_Clone(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	tree2 := tree.Clone()
	assert.Equal(t, []int{1, 2, 2, 3, 5}, tree2.ToArray())
}

func TestRBTree_ToArray(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	assert.Equal(t, []int{1, 2, 2, 3, 5}, tree.ToArray())
}

func TestRBTree_ToJSON(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	jsonBytes, err := tree.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,2,3,5]`, string(jsonBytes))
}

func TestRBTree_MarshalJSON(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	jsonBytes, err := tree.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,2,3,5]`, string(jsonBytes))
}

func TestRBTree_UnmarshalJSON(t *testing.T) {
	tree := NewRBTree(_cmp{})
	err := json.Unmarshal([]byte(`[1,2,2,3,4]`), tree)
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 2, 3, 4}, tree.ToArray())
}

func TestRBTree_String(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	str := tree.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`RBTree\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\}`, tree.Count()))
	assert.True(t, pattern.MatchString(str))
}
