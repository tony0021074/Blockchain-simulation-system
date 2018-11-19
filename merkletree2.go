package core

import (
	"crypto/sha256"
	"fmt"
)
//根节点的结构
type MerkleTree struct {
	RootNode *MerkleNode
}
//叶子节点的结构
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}
//NewMerkleNode 用来创建一棵merkle tree 的叶子节点，根据transactions
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		var data []byte
		if right != nil {
			data = append(left.Data, right.Data...)
		} else {
			data = left.Data
		}
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right
	return &mNode
}
//NewMerkleTree 创建一棵树
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	if data == nil {
		fmt.Print("args error")
	}
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	for _, idata := range data {
		node := NewMerkleNode(nil, nil, idata)
		nodes = append(nodes, *node)
	}
//遍历整棵树
	for i := 0; i < len(data)/2; i++ {
		var fathernodes []MerkleNode
		if len(nodes) == 1 {
			break
		}
		for j := 0; j < len(nodes); j += 2 {
			if j+1 >= len(nodes) {
				node := NewMerkleNode(&nodes[j], nil, nil)
				fathernodes = append(fathernodes, *node)
			} else {
				node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
				fathernodes = append(fathernodes, *node)
			}
		}

		nodes = fathernodes
	}
	mTree := MerkleTree{
		RootNode: &nodes[0],
	}
	return &mTree
}