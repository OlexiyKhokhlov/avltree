# avltree
Go maps are hash based. They have a good performance for the searching but don't provide any possibility to arange keys in the order.
Also hash-maps are hungry for the memory usage.
AVL tree's aren't so extremely fast in searching but they aren't so hungry for memory usage.
Also AVL tree's can provide a possibility to get map contens be ascending or by descending. 

## Installation
Like any other Go package just call

`$ go get github.com/OlexiyKhokhlov/avltree`

## Notes
Go is very simple language so it has any protection for dummy developer.
I can't write a code that blocks any possibility to change a key that is already inserted in AVLTree.
You must don't do that! It provides broken AVL balancing!

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
    tree := avltree.NewAVLTree(func(a interface{}, b interface{}) int {
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
    tree.Empty() //returns false
    tree.Count() //returns 2
    
    //Get values 
    tree.Contains(3) //returns false
    tree.Contains(2) //returns true
    tree.Find(2)     //returns interface{} for string that has been inserted 
    tree.Find(22)    //returns interface{} for string that has been inserted
    tree.Find(33)    //returns nil 

    //Erasing
    tree.Erase(2)
    tree.Erase(22)
}
```

## Implementations details
 - Every node has only two pointers
 - Insert and Erase methods aren't recursive

## TODO
Extend and Improve enumerating methods
