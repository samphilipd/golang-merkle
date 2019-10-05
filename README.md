# Merkle

Get the Merkle root of a file.

NOTE:

One round of SHA256 hashing is used and every node always has two children.

If, when forming a row in the tree (other than the root of the tree), it would have an odd number of elements, the final double-hash is duplicated to ensure that the row has an even number of hashes. So every node always has a left and right child unless it is a leaf in which case it has no children.

## Run the CLI

Get merkle root of the included file "alice.txt":

```
go run merkletree.go
```

Get merkle root of a specified file:

```
go run merkletree.go -file alice_trunc.txt
```

Print the hash of a node found at a particular route (0 means left, 1 means right):

```
go run merkletree.go -route 00000101
```

Help:

```
go run merkletree.go -h
```

## Run the tests

```
go test -v
```
