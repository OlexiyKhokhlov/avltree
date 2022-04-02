[![Go Reference](https://pkg.go.dev/badge/github.com/OlexiyKhokhlov/avltree.svg)](https://pkg.go.dev/github.com/OlexiyKhokhlov/avltree)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/OlexiyKhokhlov/avltree)](https://goreportcard.com/report/github.com/OlexiyKhokhlov/avltree)
# AVL Tree Go's module

## Version Information
Version 2.x.x uses generics that was introduced in the Go 1.18.
Old 1.x.x can be used in Go without generics.

## Installation
`$ go get github.com/OlexiyKhokhlov/avltree`

## Motivation
Go's standard library is still lean. Unlike other programming languages it hasn't got a container that holds a pair key-value in the determined order.
This module tries to fix it!

## Internal details
Ordered containsers usually based on [RB-Tree](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree). But it uses [AVL Tree](https://en.wikipedia.org/wiki/AVL_tree). There are a several reasons for that:
+ AVL-Tree always provides more balanced tree.
+ RB-Tree balancer is really lame in the inserting of ordered sequences.
+ RB-Tree is more popular so I don't want to implement one more of that.

This implementation has the next features:
+ Every tree node contains only two node pointers. It helps to reduce memory usage.
+ Iterators isn't implemented. May be I'll do that later. You can use enumeration methods instead.
+ Go hasn't `const` qualifier. So there is no flex way to block possibility to change a key inside of the tree. Of course I know about copying. But isn't a good solution. First of all Go hasn't got unified way to copy any type of data. And secondly it provides a bad performance when your key is a big structure. So be carefull, avoid key changing! 
+ Inserting and Erasing methods are non-recursive. So it helps to reduse stack memory usage. And also it is cool!

## Example
```
package main

import (
	"github.com/OlexiyKhokhlov/avltree"
	"strings"
	"fmt"
)

func main() {
// Create a tree where Key type is int and value type is string.
tree := avltree.NewAVLTreeOrderedKey[int, string]()

// Print true since tree has no elements yet
fmt.printLn(tree.Empty())

// Print o since tree has no elements
fmt.printLn(tree.Size())

// Insert key - 7 with value - "seven"
err := tree.Insert(7, "seven")
// err is nil 

// Insert key - 0 with value - "zero"
err = tree.Insert(0, "zero")
// err is nil 

// Insert key - 2  with value - "two"
err = tree.Insert(2, "two")
// err is nil 

// Insert key - 2  with value - "two" again
err = tree.Insert(2, "two")
// err is not nil since tree already has such key.

// Print false since tree has no 5
fmt.printLn(tree.Contains(5))

// Print true since tree has no 7
fmt.printLn(tree.Contains(7))
}
```

## TODO
+ It will be interesting to provide some performance testing for this implementation and comparing with other not only Go.

## Links
