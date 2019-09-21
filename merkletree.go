package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

const (
	path           = "alice_trunc.txt"
	chunkSizeBytes = 1024
)

// Node is a merkle tree node
// Branch nodes always have two children
// Leaf nodes have nil children
type Node struct {
	value [sha256.Size]byte
	left  *Node
	right *Node
}

func main() {
	rootNode := filePath2RootNode(path)
	fmt.Println(hex.EncodeToString(rootNode.value[:]))
}

func filePath2RootNode(path string) Node {
	f, err := os.Open(path)
	if err != nil {
		panic("Cannot open file")
	}

	fileInfo, _ := os.Stat(path)
	nChunks := nChunks(fileInfo)

	leafNodes := make([]Node, nChunks)
	buffer := make([]byte, chunkSizeBytes)

	for i := 0; i < nChunks; i++ {
		length, err := f.Read(buffer)
		if err != nil {
			panic("Unexpected error reading file")
		}
		node := makeLeafNode(buffer[:length])
		leafNodes[i] = node
	}

	return rootNode(leafNodes)
}

func nChunks(fileInfo os.FileInfo) int {
	size := fileInfo.Size()
	nChunks := int(size / chunkSizeBytes)
	if size%chunkSizeBytes != 0 {
		nChunks++
	}
	return nChunks

}

func makeLeafNode(bytes []byte) Node {
	hash := hash(bytes)
	return Node{
		hash,
		nil,
		nil,
	}
}

func rootNode(nodes []Node) Node {
	return buildTree(nodes)[0]
}

func buildTree(nodes []Node) []Node {
	if len(nodes) == 1 {
		// The root node
		return nodes
	}
	if len(nodes)%2 == 1 {
		// Double up the last node so there is always an even number
		nodes = append(nodes, nodes[len(nodes)-1])
	}
	parents := make([]Node, len(nodes)/2)
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		right := nodes[i+1]
		parents[i/2] = makeParent(left, right)
	}
	return buildTree(parents)
}

func makeParent(left, right Node) Node {
	var concatenatedHashes [2 * sha256.Size]byte
	copy(concatenatedHashes[:], left.value[:])
	copy(concatenatedHashes[sha256.Size:], right.value[:])
	hash := hash(concatenatedHashes[:])
	return Node{
		hash,
		&left,
		&right,
	}
}

func hash(block []byte) [sha256.Size]byte {
	return sha256.Sum256(block)
}
