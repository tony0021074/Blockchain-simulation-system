package main

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
)

//MerkleTreeRoot : Calculating Merkle Tree Root
func MerkleTreeRoot(content []string) (headRoot []byte) {
	var mkt MerkleTree
	return mkt.GenerateRoot(content).NodeHash
}

// MerkleNode : Node of Merkle Tree
type MerkleNode struct {
	LeftNode  *MerkleNode
	RightNode *MerkleNode
	NodeData  string
	NodeHash  []byte
}

// MerkleTree : Root Node
type MerkleTree struct {
	Level    int
	RootNode *MerkleNode
}

// CalLevel : Calculate the deep of Merkle Tree
func (mkt MerkleTree) CalLevel(data []string) (level int) {
	level = 0
	for {
		if len(data) > int(math.Pow(2, float64(level))) {
			level = level + 1
		} else {
			break
		}
	}
	return level
}

// CalSHA256Hash : Calculate a sha256 hash
func (mkt MerkleTree) CalSHA256Hash(input string) []byte {
	h := sha256.New()
	h.Write([]byte(input))
	return h.Sum(nil)
}

// GenerateRoot : Create Merkle Tree
func (mkt MerkleTree) GenerateRoot(data []string) *MerkleNode {

	// Prepare Leaf Node
	mkt.Level = mkt.CalLevel(data)
	var nodes []*MerkleNode
	for i := 0; i < len(data); i++ {
		nodes = append(nodes, &MerkleNode{
			NodeData: data[i],
			NodeHash: mkt.CalSHA256Hash(data[i]),
		})
	}
	for i := len(data); i < int(math.Pow(2, float64(mkt.Level))); i++ {
		nodes = append(nodes, &MerkleNode{
			NodeData: data[len(data)-1],
			NodeHash: mkt.CalSHA256Hash(data[len(data)-1]),
		})
	}
	//Building Tree
	for {
		if len(nodes) != 1 {
			var tempNodes []*MerkleNode
			for i := 0; i < len(nodes); i = i + 2 {
				tempNodes = append(tempNodes, &MerkleNode{
					LeftNode:  nodes[i],
					RightNode: nodes[i+1],
					NodeData:  hex.EncodeToString(nodes[i].NodeHash) + hex.EncodeToString(nodes[i+1].NodeHash),
					NodeHash:  mkt.CalSHA256Hash(hex.EncodeToString(nodes[i].NodeHash) + hex.EncodeToString(nodes[i+1].NodeHash)),
				})
			}
			nodes = tempNodes
		} else {
			break
		}
	}
	mkt.RootNode = nodes[0]
	return mkt.RootNode
}
