# Merkle

Get the Merkle root of the file at ./alice.txt.

## Run the CLI

Get merkle root of the file:

```
go run merkletree.go
```

Print the hash of a node found at a particular route (0 means left, 1 means right):

```
go run merkletree.go 00000101
```

NOTE: If, when forming a row in the tree (other than the root of the tree), it would have an odd number of elements, the final double-hash is duplicated to ensure that the row has an even number of hashes. So every node always has a left and right child unless it is a leaf in which case it has no children.

## Run the tests

```
go test -v
```
