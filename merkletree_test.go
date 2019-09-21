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
