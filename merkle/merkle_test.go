package merkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerkleTree(t *testing.T) {
	// 创建一些测试数据
	data := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("example"),
		[]byte("test"),
		[]byte("hellofe"),
		[]byte("worfdwfld"),
		[]byte("exfeample"),
		[]byte("tfeest"),
	}

	// 创建Merkle树
	merkleTree := NewMerkelTree(data)

	// 验证Merkle树的每个节点
	for i := 0; i < len(data); i++ {
		assert.True(t, merkleTree.Validate(i, data[i]))
	}

	// 验证一个错误的情况
	invalidData := []byte("world")
	assert.False(t, merkleTree.Validate(0, invalidData))
}
