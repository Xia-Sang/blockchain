// 实现 Merkle 树
package merkle

import (
	"bytes"
	"case/utils"
)

// node节点的细节[所有细节不对外暴露]
type node struct {
	data                []byte
	hash                []byte
	left, right, parent *node
}

// newNode 新建节点
func newNode(data []byte) *node {
	return &node{data: data, hash: utils.GetDataHash(data)}
}

// NewMerkelTree [不对外暴露]
type MerkelTree struct {
	nodes []*node
	root  *node
}

// 新建MerkelTree
func NewMerkelTree(data [][]byte) *MerkelTree {
	var nodes []*node
	for _, d := range data {
		nodes = append(nodes, newNode(d))
	}
	return &MerkelTree{nodes: nodes}
}

// 实现建立工作【不对外暴露细节】
func (mk *MerkelTree) buildTree() {
	if len(mk.nodes) == 0 {
		mk.root = nil
		return
	}
	queue := mk.nodes
	for len(queue) > 1 {
		left := queue[0]
		right := queue[1]
		queue = queue[2:]
		newnode := newNode(append(left.hash, right.hash...))
		newnode.left = left
		newnode.right = right
		left.parent = newnode
		right.parent = newnode
		queue = append(queue, newnode)
	}
	mk.root = queue[0]
}

// 获取根节点
func (mk *MerkelTree) getMerkelRoot() *node {
	if mk.root == nil && len(mk.nodes) > 0 {
		mk.buildTree()
	}
	return mk.root
}

// 计数node节点【叶子节点】
func (mk *MerkelTree) countMerkelNodes() int {
	if mk.root == nil && len(mk.nodes) > 0 {
		mk.buildTree()
	}
	return len(mk.nodes)
}

// 对于只暴露验证节点
func (mk *MerkelTree) Validate(index int, data []byte) bool {
	if mk.nodes == nil || len(mk.nodes) == 0 {
		mk.buildTree()
	}
	if index < 0 || index >= mk.countMerkelNodes() {
		return false
	}
	curNode := mk.nodes[index]
	if !bytes.Equal(utils.GetDataHash([]byte(data)), curNode.hash) {
		return false
	}
	for curNode != nil && curNode.parent != nil {
		par := curNode.parent
		var newData []byte
		if par.left == curNode {
			newData = append(curNode.hash, par.right.hash...)
		}
		if par.right == curNode {
			newData = append(par.left.hash, curNode.hash...)
		}
		if !bytes.Equal(utils.GetDataHash(newData), par.hash) {
			return false
		}
		curNode = par
	}
	return bytes.Equal(curNode.hash, mk.getMerkelRoot().hash)
}
