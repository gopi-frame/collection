package tree

import "github.com/gopi-frame/contract"

const (
	red   = true
	black = false
)

type rbNode[E any] struct {
	value E
	left  *rbNode[E]
	right *rbNode[E]
	color bool
	count int
}

func (node *rbNode[E]) leftRotate() *rbNode[E] {
	if node == nil {
		return nil
	}
	pivot := node.right
	node.right = pivot.left
	pivot.left = node
	pivot.color = node.color
	node.color = red
	return pivot
}

func (node *rbNode[E]) rightRotate() *rbNode[E] {
	if node == nil {
		return nil
	}
	pivot := node.left
	node.left = pivot.right
	pivot.right = node
	pivot.color = node.color
	node.color = red
	return pivot
}

func (node *rbNode[E]) isRed() bool {
	if node == nil {
		return false
	}
	return node.color
}

func (node *rbNode[E]) isBlack() bool {
	if node == nil {
		return true
	}
	return !node.color
}

func (node *rbNode[E]) switchColor() {
	if node == nil {
		return
	}
	node.color = !node.color
	node.left.color = !node.left.color
	node.right.color = !node.right.color
}

func (node *rbNode[E]) moveRedLeft() *rbNode[E] {
	node.switchColor()
	if node.right.left.isRed() {
		node.right = node.right.rightRotate()
		node = node.leftRotate()
		node.switchColor()
	}
	return node
}

func (node *rbNode[E]) moveRedRight() *rbNode[E] {
	node.switchColor()
	if node.left.left.isRed() {
		node = node.rightRotate()
		node.switchColor()
	}
	return node
}

func (node *rbNode[E]) insert(value E, comparator contract.Comparator[E]) *rbNode[E] {
	if node == nil {
		return &rbNode[E]{
			value: value,
			color: red,
			count: 1,
		}
	}
	result := comparator.Compare(value, node.value)
	if result == 0 {
		node.count++
		return node
	} else if result < 0 {
		node.left = node.left.insert(value, comparator)
	} else {
		node.right = node.right.insert(value, comparator)
	}
	activeNode := node
	if activeNode.right.isRed() && activeNode.left.isBlack() {
		activeNode = node.leftRotate()
	} else {
		if activeNode.left.isRed() && activeNode.left.left.isRed() {
			activeNode = activeNode.rightRotate()
		}
		if activeNode.left.isRed() && activeNode.right.isRed() {
			activeNode.switchColor()
		}
	}
	return activeNode
}

func (node *rbNode[E]) remove(value E, comparator contract.Comparator[E]) *rbNode[E] {
	activeNode := node
	if comparator.Compare(value, node.value) < 0 {
		if activeNode.left.isBlack() && activeNode.left.left.isBlack() {
			activeNode = activeNode.moveRedLeft()
		}
		activeNode.left = activeNode.left.remove(value, comparator)
	} else {
		if activeNode.left.isRed() {
			activeNode = activeNode.rightRotate()
		}
		if comparator.Compare(value, activeNode.value) == 0 && activeNode.right == nil {
			return nil
		}
		if activeNode.right.isBlack() && activeNode.right.left.isBlack() {
			activeNode = activeNode.moveRedRight()
		}
		if comparator.Compare(value, activeNode.value) == 0 {
			min := activeNode.right.min()
			activeNode.value = min.value
			activeNode.count = min.count
			activeNode.right = activeNode.right.removeMin()
		} else {
			activeNode.right = activeNode.right.remove(value, comparator)
		}
	}
	return activeNode.fix()
}

func (node *rbNode[E]) removeMin() *rbNode[E] {
	activeNode := node
	if activeNode.left == nil {
		return nil
	}
	if activeNode.left.isBlack() && activeNode.left.left.isBlack() {
		activeNode = activeNode.moveRedLeft()
	}
	activeNode.left = activeNode.left.removeMin()
	return activeNode.fix()
}

func (node *rbNode[E]) fix() *rbNode[E] {
	activeNode := node
	if activeNode.right.isRed() {
		activeNode = activeNode.leftRotate()
	}
	if activeNode.left.isRed() && activeNode.left.left.isRed() {
		activeNode = activeNode.rightRotate()
	}
	if activeNode.left.isRed() && activeNode.right.isRed() {
		activeNode.switchColor()
	}
	return activeNode
}

func (node *rbNode[E]) min() *rbNode[E] {
	if node.left == nil {
		return node
	}
	return node.left.min()
}

func (node *rbNode[E]) max() *rbNode[E] {
	if node.right == nil {
		return node
	}
	return node.right.max()
}

func (node *rbNode[E]) find(value E, comparator contract.Comparator[E]) *rbNode[E] {
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

func (node *rbNode[E]) inOrderRange() (nodes []*rbNode[E]) {
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
