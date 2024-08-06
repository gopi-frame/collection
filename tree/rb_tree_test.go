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
	tree := NewRBTree[int](_cmp{})
	assert.True(t, tree.IsEmpty())
}

func TestRBTree_IsNotEmpty(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 4)
	assert.True(t, tree.IsNotEmpty())
}

func TestRBTree_Contains(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		ok := tree.Contains(1)
		assert.False(t, ok)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewRBTree(_cmp{}, 1, 2, 3)
		ok := tree.Contains(1)
		assert.True(t, ok)
	})
}

func TestRBTree_Remove(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		tree.Remove(1)
		assert.Equal(t, int64(0), tree.Count())
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewRBTree(_cmp{}, 1, 2, 3)
		tree.Remove(1)
		assert.Equal(t, int64(2), tree.Count())
		assert.False(t, tree.Contains(1))
	})
}

func TestRBTree_Clear(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3)
	assert.False(t, tree.IsEmpty())
	tree.Clear()
	assert.True(t, tree.IsEmpty())
}

func TestRBTree_First(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		v, ok := tree.First()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewRBTree(_cmp{}, 1, 2, 3)
		v, ok := tree.First()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
	})
}

func TestRBTree_FirstOr(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		v := tree.FirstOr(1)
		assert.Equal(t, 1, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewRBTree(_cmp{}, 2, 3)
		v := tree.FirstOr(1)
		assert.Equal(t, 2, v)
	})
}

func TestRBTree_Last(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		v, ok := tree.Last()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewRBTree(_cmp{}, 1, 2, 3)
		v, ok := tree.Last()
		assert.True(t, ok)
		assert.Equal(t, 3, v)
	})
}

func TestRBTree_LastOr(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		v := tree.LastOr(1)
		assert.Equal(t, 1, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewRBTree(_cmp{}, 2, 3)
		v := tree.LastOr(10)
		assert.Equal(t, 3, v)
	})
}

func TestRBTree_Each(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	var items []int
	tree.Each(func(_ int, value int) bool {
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
	t.Run("valid json", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		err := json.Unmarshal([]byte(`[1,2,2,3,4]`), tree)
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2, 2, 3, 4}, tree.ToArray())
	})

	t.Run("invalid json", func(t *testing.T) {
		tree := NewRBTree[int](_cmp{})
		err := json.Unmarshal([]byte(`[1,2,2,3,4`), tree)
		assert.NotNil(t, err)
	})
}

func TestRBTree_String(t *testing.T) {
	tree := NewRBTree(_cmp{}, 1, 2, 3, 5, 2)
	str := tree.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`RBTree\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\}`, tree.Count()))
	assert.True(t, pattern.MatchString(str))
}
