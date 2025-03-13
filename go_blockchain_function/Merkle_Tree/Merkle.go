package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

type MerkleTree struct {
	RootNode *MerkleNode
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right
	return &mNode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}
	mTree := MerkleTree{&nodes[0]}

	return &mTree
}

func showMerkleTree(root *MerkleNode) {
	if root == nil {
		return
	} else {
		PrintNode(root)
	}
	showMerkleTree(root.Left)
	showMerkleTree(root.Right)
}

func check(node *MerkleNode) bool {
	if node.Left == nil {
		return true
	}
	prevHashes := append(node.Left.Data, node.Right.Data...)
	hash32 := sha256.Sum256(prevHashes)
	hash := hash32[:]
	return bytes.Compare(hash, node.Data) == 0
}

func PrintNode(node *MerkleNode) {
	fmt.Printf("%p\n", node)
	if node != nil {
		fmt.Printf("left[%p],right[%p],data(%x)\n", node.Left, node.Right, node.Data)
		fmt.Printf("check:%v\n", check(node))
	}
}

func main() {
	data := [][]byte{
		[]byte("1000"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
	}
	mTree := NewMerkleTree(data)
	showMerkleTree(mTree.RootNode)
}
