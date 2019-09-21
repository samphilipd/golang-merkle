package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

const (
	// Path is path to the file for which we want to get the Merkle root
	Path = "alice.txt"
	// ChunkSize is chunk size in bytes
	ChunkSize = 1024
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
	var lookupRoute string

	if len(os.Args) == 2 {
		lookupRoute = os.Args[1]
	} else if len(os.Args) != 1 {
		panic("Only 0 or 1 argument is allowed")
	}

	rootNode := filePath2RootNode(Path)

	node := lookup(rootNode, lookupRoute)

	fmt.Println(hex.EncodeToString(node.value[:]))
}

func lookup(rootNode Node, route string) Node {
	node := rootNode

	for i := 0; i < len(route); i++ {
		if node.left == nil {
			depth := depth(rootNode)
			panic(fmt.Sprintf("You tried to lookup %v nodes deep, the max depth of this tree is %v", len(route), depth))
		}
		if route[i] == '0' {
			node = *node.left
		} else if route[i] == '1' {
			node = *node.right
		} else {
			panic("Only 0 or 1 is allowed in routes")
		}
	}
	return node
}

func depth(node Node) int {
	return dive(node, 0)
}

func dive(node Node, currentDepth int) int {
	if node.left != nil {
		return dive(*node.left, currentDepth+1)
	}
	return currentDepth
}

func filePath2RootNode(path string) Node {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fileInfo, _ := f.Stat()
	nChunks := nChunks(fileInfo)

	leafNodes := make([]Node, nChunks)
	buffer := make([]byte, ChunkSize)

	for i := 0; i < nChunks; i++ {
		length, err := f.Read(buffer)
		if err != nil {
			panic(err)
		}
		node := makeLeaf(buffer[:length])
		leafNodes[i] = node
	}

	return rootNode(leafNodes)
}

func nChunks(fileInfo os.FileInfo) int {
	size := fileInfo.Size()
	nChunks := int(size / ChunkSize)
	if size%ChunkSize != 0 {
		nChunks++
	}
	return nChunks

}

func makeLeaf(bytes []byte) Node {
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
		// We always eventually return an array with one node in it - the root node
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
