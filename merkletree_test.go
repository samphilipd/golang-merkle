package main

import (
	"encoding/hex"
	"testing"
)

func TestFilePath2RootNode(t *testing.T) {
	rootNode := filePath2RootNode("alice.txt")
	expected := "5a792e5fab93b0bbbdee42f10e40b2a05296cd301144620bee33999c8115eefb"
	actual := hex.EncodeToString(rootNode.value[:])
	if expected != actual {
		t.Fatalf("filePath2RootNode failed: expect %v got %v", expected, actual)
	}
}

func TestLookup(t *testing.T) {
	rootNode := filePath2RootNode("alice.txt")
	node := lookup(rootNode, "00000101")
	expected := "a144cb1e6029f92b0403071a07b69dd9d3aff267c8ff5f2e6e950721562bf84b"
	actual := hex.EncodeToString(node.value[:])
	if expected != actual {
		t.Fatalf("lookup failed: expect %v got %v", expected, actual)
	}
}
