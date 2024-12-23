package tree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type _cmp struct{}

func (c _cmp) Compare(a, b int) int {
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func TestAVLTree_Count(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3)
	assert.Equal(t, int64(3), tree.Count())
}

func TestAVLTree_IsEmpty(t *testing.T) {
	tree := NewAVLTree[int](_cmp{})
	assert.True(t, tree.IsEmpty())
}

func TestAVLTree_IsNotEmpty(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 4)
	assert.True(t, tree.IsNotEmpty())
}

func TestAVLTree_Contains(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		ok := tree.Contains(1)
		assert.False(t, ok)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewAVLTree(_cmp{}, 1, 2, 3)
		ok := tree.Contains(1)
		assert.True(t, ok)
	})
}

func TestAVLTree_Remove(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		tree.Remove(1)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewAVLTree(_cmp{}, 1, 2, 3)
		tree.Remove(1)
		assert.Equal(t, int64(2), tree.Count())
		assert.False(t, tree.Contains(1))
	})
}

func TestAVLTree_Clear(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3)
	assert.False(t, tree.IsEmpty())
	tree.Clear()
	assert.True(t, tree.IsEmpty())
}

func TestAVLTree_First(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		v, ok := tree.First()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewAVLTree(_cmp{}, 1, 2, 3)
		v, ok := tree.First()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
	})
}

func TestAVLTree_FirstOr(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		v := tree.FirstOr(1)
		assert.Equal(t, 1, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewAVLTree(_cmp{}, 2, 3)
		v := tree.FirstOr(1)
		assert.Equal(t, 2, v)
	})
}

func TestAVLTree_Last(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		v, ok := tree.Last()
		assert.False(t, ok)
		assert.Equal(t, 0, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewAVLTree(_cmp{}, 1, 2, 4, 2)
		v, ok := tree.Last()
		assert.True(t, ok)
		assert.Equal(t, 4, v)
	})
}

func TestAVLTree_LastOr(t *testing.T) {
	t.Run("empty tree", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		v := tree.LastOr(1)
		assert.Equal(t, 1, v)
	})

	t.Run("non-empty tree", func(t *testing.T) {
		tree := NewAVLTree(_cmp{}, 1, 2, 4, 2)
		v := tree.LastOr(10)
		assert.Equal(t, 4, v)
	})
}

func TestAVLTree_Each(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 5, 2)
	var items []int
	tree.Each(func(_ int, value int) bool {
		items = append(items, value)
		return value < 2
	})
	assert.Equal(t, []int{1, 2}, items)
}

func TestAVLTree_Clone(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 5, 2)
	tree2 := tree.Clone()
	assert.Equal(t, []int{1, 2, 2, 3, 5}, tree2.ToArray())
}

func TestAVLTree_ToArray(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 5, 2)
	assert.Equal(t, []int{1, 2, 2, 3, 5}, tree.ToArray())
}

func TestAVLTree_ToJSON(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 5, 2)
	jsonBytes, err := tree.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,2,3,5]`, string(jsonBytes))
}

func TestAVLTree_MarshalJSON(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 5, 2)
	jsonBytes, err := tree.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,2,3,5]`, string(jsonBytes))
}

func TestAVLTree_UnmarshalJSON(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		err := json.Unmarshal([]byte(`[1,2,2,3,4]`), tree)
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2, 2, 3, 4}, tree.ToArray())
	})

	t.Run("invalid json", func(t *testing.T) {
		tree := NewAVLTree[int](_cmp{})
		err := json.Unmarshal([]byte(`[1,2,2,3,4`), tree)
		assert.NotNil(t, err)
	})
}

func TestAVLTree_String(t *testing.T) {
	tree := NewAVLTree(_cmp{}, 1, 2, 3, 5, 2)
	str := tree.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`AVLTree\[int\]\(len=%d\)\{\n(\t\d+,\n){5}\}`, tree.Count()))
	assert.True(t, pattern.MatchString(str))
}
