# avltree
This is AVL Tree implementations.

## Installation
`$ go get github.com/OlexiyKhokhlov/avltree`

## How to use
```
package main

import (
    "fmt"
    "github.com/OlexiyKhokhlov/avltree"
)

func main() {
    //Tree creation.
    //You have to pass into NewAVLTree a Key comparator function
    //Example uses int comparator. So Key type is int here
    tree := NewAVLTree(func(a interface{}, b interface{}) int {
	first := a.(int)
	second := b.(int)

	if first == second {
	    return 0
	}
	if first < second {
	    return -1
	}
	return 1
    })
    
    //Inserting
    tree.Insert(2, "this is value for key 2")
    tree.Insert(22, "this is value for key 22")
    
    //Tree info
    tree.Empty()
    tree.Count()

    //Erasing
    tree.Erase(2)
    tree.Erase(22)
}
```

## Implementations details
 - Every node has only two pinters
 - Insert and Erase methods aren't recursive

## TODO
Order enumerating methods
