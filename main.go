package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const (
	path           = "alice.txt"
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
	f, err := os.Open(path)
	if err != nil {
		panic("Cannot open file")
	}

	// TODO: Make the slice have size the length of the file
	leafNodes := make([]Node, 0)

	for {
		buffer := make([]byte, chunkSizeBytes)
		length, err := f.Read(buffer)
		if err == io.EOF {
			break
		}
		// fmt.Println(buffer)
		// fmt.Printf("%d bytes\n", length)
		node := makeLeafNode(buffer[:length])
		leafNodes = append(leafNodes, node)
	}

	rootNode := rootNode(leafNodes)
	fmt.Println(hex.EncodeToString(rootNode.value[:]))
}

func makeLeafNode(bytes []byte) Node {
	hash := hash(bytes)
	// fmt.Println(hex.EncodeToString(hash[:]))
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
	// fmt.Println(len(nodes))
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
	hash := sha256.Sum256(block)
	return hash
}

func printValue(node Node) {
	fmt.Println(node.value)
}
