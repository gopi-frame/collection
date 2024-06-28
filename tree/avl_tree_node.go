package tree

import (
	"github.com/gopi-frame/contract"
)

type avlNode[E any] struct {
	value  E
	left   *avlNode[E]
	right  *avlNode[E]
	height int
	count  int
}

func (node *avlNode[E]) updateHeight() {
	var leftHeight, rightHeight = 0, 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	max := leftHeight
	if rightHeight > leftHeight {
		max = rightHeight
	}
	node.height = max + 1
}

func (node *avlNode[E]) drop() int {
	var leftHeight, rightHeight = 0, 0
	if node.left != nil {
		leftHeight = node.left.height
	}
	if node.right != nil {
		rightHeight = node.right.height
	}
	return leftHeight - rightHeight
}

func (node *avlNode[E]) insert(value E, comparator contract.Comparator[E]) *avlNode[E] {
	if node == nil {
		return &avlNode[E]{
			value:  value,
			height: 1,
			count:  1,
		}
	}
	if comparator.Compare(value, node.value) == 0 {
		node.count++
		return node
	}
	var newNode *avlNode[E]
	if comparator.Compare(value, node.value) < 0 {
		node.left = node.left.insert(value, comparator)
		if node.drop() == 2 {
			if comparator.Compare(value, node.left.value) < 0 {
				newNode = node.rightRotate()
			} else {
				newNode = node.leftRightRotate()
			}
		}
	} else {
		node.right = node.right.insert(value, comparator)
		if node.drop() == -2 {
			if comparator.Compare(value, node.right.value) < 0 {
				newNode = node.rightLeftRotate()
			} else {
				newNode = node.leftRotate()
			}
		}
	}
	if newNode == nil {
		node.updateHeight()
		return node
	}
	newNode.updateHeight()
	return newNode
}

func (node *avlNode[E]) leftRotate() *avlNode[E] {
	pivot := node.right
	node.right = pivot.left
	pivot.left = node
	node.updateHeight()
	pivot.updateHeight()
	return pivot
}

func (node *avlNode[E]) rightRotate() *avlNode[E] {
	pivot := node.left
	node.left = pivot.right
	pivot.right = node
	node.updateHeight()
	pivot.updateHeight()
	return pivot
}

func (node *avlNode[E]) leftRightRotate() *avlNode[E] {
	node.left = node.left.leftRotate()
	return node.rightRotate()
}

func (node *avlNode[E]) rightLeftRotate() *avlNode[E] {
	node.right = node.right.rightRotate()
	return node.leftRotate()
}

func (node *avlNode[E]) find(value E, comparator contract.Comparator[E]) *avlNode[E] {
	if node == nil {
		return nil
	}
	result := comparator.Compare(value, node.value)
	if result == 0 {
		return node
	} else if result < 0 {
		return node.left.find(value, comparator)
	} else {
		return node.right.find(value, comparator)
	}
}

func (node *avlNode[E]) min() *avlNode[E] {
	if node.left == nil {
		return node
	}
	return node.left.min()
}

func (node *avlNode[E]) max() *avlNode[E] {
	if node.right == nil {
		return node
	}
	return node.right.max()
}

func (node *avlNode[E]) remove(value E, comparator contract.Comparator[E]) *avlNode[E] {
	if node == nil {
		return nil
	}
	result := comparator.Compare(value, node.value)
	if result < 0 {
		node.left = node.left.remove(value, comparator)
	} else if result > 0 {
		node.right = node.right.remove(value, comparator)
	} else {
		if node.left == nil && node.right == nil {
			return nil
		}
		if node.left != nil && node.right != nil {
			if node.left.height > node.right.height {
				max := node.left.max()
				node.value = max.value
				node.count = max.count
				node.left = node.left.remove(max.value, comparator)
			} else {
				min := node.right.min()
				node.value = min.value
				node.count = min.count
				node.right = node.right.remove(min.value, comparator)
			}
		} else if node.left != nil {
			node.value = node.left.value
			node.count = node.left.count
			node.height = 1
			node.left = nil
		} else {
			node.value = node.right.value
			node.count = node.right.count
			node.height = 1
			node.right = nil
		}
		return node
	}
	var newNode *avlNode[E]
	drop := node.drop()
	if drop == 2 {
		if node.left.drop() == 1 {
			newNode = node.rightRotate()
		} else {
			newNode = node.leftRightRotate()
		}
	} else if drop == -2 {
		if node.right.drop() == -1 {
			newNode = node.leftRotate()
		} else {
			newNode = node.rightLeftRotate()
		}
	}
	if newNode == nil {
		node.updateHeight()
		return node
	}
	newNode.updateHeight()
	return newNode
}

func (node *avlNode[E]) inOrderRange() (nodes []*avlNode[E]) {
	if node == nil {
		return
	}
	nodes = append(nodes, node.left.inOrderRange()...)
	for i := 0; i < node.count; i++ {
		nodes = append(nodes, node)
	}
	nodes = append(nodes, node.right.inOrderRange()...)
	return
}
