[![Go Reference](https://pkg.go.dev/badge/gopkg.in/OlexiyKhokhlov/avltree.v2.svg)](https://pkg.go.dev/gopkg.in/OlexiyKhokhlov/avltree.v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/OlexiyKhokhlov/avltree.v2)](https://goreportcard.com/report/gopkg.in/OlexiyKhokhlov/avltree.v2)
[![Go version](https://img.shields.io/badge/go-v1.18-blue)](https://golang.org/dl/#stable)
# AVL Tree Go's module

## Version Information
The current version 2.0.0 has the same functionality as v1.0.4 but uses genrecis that was introduced in Go v1.18.

Old version  1.0.4 is placed here `https://github.com/OlexiyKhokhlov/avltree/tree/main`

## Installation
`$ go get gopkg.in/OlexiyKhokhlov/avltree.v2`

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
+ Go hasn't `const` qualifier. So there is no flex way to block possibility to key changinging inside of the tree. Of course I know about copying. But isn't a good solution. First of all Go hasn't got unified way to copy any type of data. And secondly it provides a bad performance when your key is a big structure. So be carefull, avoid key changing! 
+ Inserting and Erasing methods are non-recursive. So it helps to reduse stack memory usage. And also it is cool!

## Examples

<details>
  <summary>Tree construction</summary>

```Go

package main

import (
	"fmt"
	"strings"

	"gopkg.in/OlexiyKhokhlov/avltree.v2"
)

func main() {
	// Create a tree where Key type is 'int' and value type is 'string'.
	// It is possible to use NewAVLTreeOrderedKey since `int` type is one from `constraints.Ordered`
	tree1 := avltree.NewAVLTreeOrderedKey[int, string]()

	// Create a tree where Key type is `*int` and value type is 'string'.
	// It is possible to use NewAVLTreeOrderedKeyPtr since `int` type is one from `constraints.Ordered`
	tree2 := avltree.NewAVLTreeOrderedKeyPtr[int, string]()

	// Create a tree where Key type is a custom struct 'MyStruct1'
	// Is much better to use Key like '*MyStruct1' since it helps to avoid internal keys copying.
	// Need to use NewAVLTree with user defined comparator for that
	type MyStruct1 struct {
		key     int
		payload string
	}
	tree3 := avltree.NewAVLTree[*MyStruct1, int](func(a *MyStruct1, b *MyStruct1) int {
		if a.key == b.key {
			return 0
		}
		if a.key < b.key {
			return -1
		}
		return 1
	})

	// Create a tree where key is a struct and a key is divided on two parts
	type MyStruct2 struct {
		KeyPart1 string
		KeyPart2 int
		Payload  string
	}

	// It is much better to use pointer when key isn't a trivial data type.
	// It helps to avoid internal copying.
	tree4 := avltree.NewAVLTree[*MyStruct2, string](func(a *MyStruct2, b *MyStruct2) int {
		strcmp := strings.Compare(a.KeyPart1, b.KeyPart1)
		if strcmp != 0 {
			return strcmp
		}
		if a.KeyPart2 == b.KeyPart2 {
			return 0
		}
		if a.KeyPart1 < b.KeyPart1 {
			return -1
		}
		return 1
	})

	//Cheating Go about not used variables
	fmt.Println(tree1.Empty())
	fmt.Println(tree2.Empty())
	fmt.Println(tree3.Empty())
	fmt.Println(tree4.Empty())
}
```
</details>

<details>
  <summary>Modification</summary>

```Go
package main

import (
	"fmt"

	"gopkg.in/OlexiyKhokhlov/avltree.v2"
)

func main() {
	// Create a tree where Key type is 'int' and value type is 'string'.
	// It is possible to use NewAVLTreeOrderedKey since `int` type is one from `constraints.Ordered`
	tree := avltree.NewAVLTreeOrderedKey[int, string]()

	//Insert key, value pairs
	err := tree.Insert(10, "10")
	// err is nil here
	if err != nil {
		fmt.Println("Insertion failed: ", err)
	} else {
		fmt.Println("Inserted")
	}

	err = tree.Insert(0, "0")
	// err is nil here
	if err != nil {
		fmt.Println("Insertion failed: ", err)
	} else {
		fmt.Println("Inserted")
	}

	err = tree.Insert(5, "5")
	// err is nil here
	if err != nil {
		fmt.Println("Insertion failed: ", err)
	} else {
		fmt.Println("Inserted")
	}

	//Now tree contains {{0,"0"}, {5,"5"}, {10, "10"}}

	//Try insert duplicate
	err = tree.Insert(5, "5")
	// err is not nil here
	if err != nil {
		fmt.Println("Insertion failed: ", err)
	} else {
		fmt.Println("Inserted")
	}

	// Erase elements by key
	err = tree.Erase(0)
	// err is nil here
	if err != nil {
		fmt.Println("Erasing failed: ", err)
	} else {
		fmt.Println("Erased")
	}

	//Now tree doesn't contains '0' key. Try remove it again
	err = tree.Erase(0)
	// err is not nil here
	if err != nil {
		fmt.Println("Erasing failed: ", err)
	} else {
		fmt.Println("Erased")
	}

	// Is possible to modificate a value that is stored inside a tree
	value := tree.Find(5)
	if value != nil {
		// if tree has a value change it
		*value = "new value"
	}
	// Find it again and print
	value = tree.Find(5)
	if value != nil {
		fmt.Println("Changed value is: ", *value)
	}

	// Clear entire a tree
	tree.Clear()
}
```
</details>

<details>
  <summary>Element access</summary>

```Go
package main

import (
	"fmt"
	"strconv"

	"gopkg.in/OlexiyKhokhlov/avltree.v2"
)

func main() {
	// Create a tree where Key type is 'int' and value type is 'string'.
	// It is possible to use NewAVLTreeOrderedKey since `int` type is one from `constraints.Ordered`
	tree := avltree.NewAVLTreeOrderedKey[int, string]()
	// Insert a set of keys
	for i := 0; i <= 100; i += 10 {
		tree.Insert(i, strconv.Itoa(i))
	}

	// Now tree contains [0, 10, ... 100] keys with corresponded values

	// Check if a tree contains some key
	if tree.Contains(5) {
		fmt.Println("Tree contains 5")
	} else {
		fmt.Println("Tree hasn't 5")
	}

	// Get pointer on stored in the tree value by the given key
	val := tree.Find(50)
	// If tree hasn't got such key returns nill
	if val != nil {
		fmt.Println("Tree contains value for key=50: ", *val)
	}

	// Get first element in the tree
	k, v := tree.First()
	// If tree is empty can return nil, nil
	if k != nil {
		fmt.Println("First element is: ", *k, *v)
	}

	// Get last element in the tree
	k, v = tree.Last()
	// If tree is empty can return nil, nil
	if k != nil {
		fmt.Println("First element is: ", *k, *v)
	}

	// Get a next element after the given key
	// The given key not necessary has been stored in the tree
	k, v = tree.FindNextElement(1)
	// If the given key is the last can return nil, nil
	if k != nil {
		fmt.Println("Element after 1 is: ", *k, *v)
	}

	// Get a prev element after the given key
	// The given key not necessary has been stored in the tree
	k, v = tree.FindPrevElement(100)
	// If the given key is the last can return nil, nil
	if k != nil {
		fmt.Println("Element before 100 is: ", *k, *v)
	}
}
```
</details>

<details>
  <summary>Enumeration</summary>

```Go
package main

import (
	"fmt"
	"strconv"

	"gopkg.in/OlexiyKhokhlov/avltree.v2"
)

func main() {
	// Create a tree where Key type is 'int' and value type is 'string'.
	// It is possible to use NewAVLTreeOrderedKey since `int` type is one from `constraints.Ordered`
	tree := avltree.NewAVLTreeOrderedKey[int, string]()

	// Insert a set of keys
	for i := 0; i <= 100; i += 10 {
		tree.Insert(i, strconv.Itoa(i))
	}
	// Now tree contains [0, 10, ... 100] keys with corresponded values

	// Print all elements in the ascending order
	fmt.Println("Ascending:")
	tree.Enumerate(avltree.ASCENDING, func(key int, value string) bool {
		fmt.Println("Element: ", key, " ", value)
		return true // Always return true since we don't want to interupt enumearation
	})

	// Print all elements in the descending order
	fmt.Println("Descending:")
	tree.Enumerate(avltree.ASCENDING, func(key int, value string) bool {
		fmt.Println("Element: ", key, " ", value)
		return true // Always return true since we don't want to interupt enumearation
	})

	// Print all element these are between start and finish in ascending order
	start := 20
	finish := 68
	fmt.Println("Diapason: 20..68:")
	tree.EnumerateDiapason(&start, &finish, avltree.ASCENDING, func(key int, value string) bool {
		fmt.Println("Element: ", key, " ", value)
		return true // Always return true since we don't want to interupt enumearation
	})

	// Print all element these are greater than start in ascending order
	start = 20
	fmt.Println("Diapason: 20...:")
	tree.EnumerateDiapason(&start, nil, avltree.ASCENDING, func(key int, value string) bool {
		fmt.Println("Element: ", key, " ", value)
		return true // Always return true since we don't want to interupt enumearation
	})

	// Print all 3 element these are greater than start in ascending order
	start = 55
	fmt.Println("First 3 in diapason: 55...:")
	i := 0
	tree.EnumerateDiapason(&start, nil, avltree.ASCENDING, func(key int, value string) bool {
		fmt.Println("Element: ", key, " ", value)
		i += 1
		if i == 3 {
			return false // 3 element is already printed. Return false for stop
		}
		return true
	})

}
```
</details>

<details>
  <summary>Capacity</summary>

```Go
package main

import (
	"fmt"
	"strconv"

	"gopkg.in/OlexiyKhokhlov/avltree.v2"
)

func main() {
	// Create a tree where Key type is 'int' and value type is 'string'.
	// It is possible to use NewAVLTreeOrderedKey since `int` type is one from `constraints.Ordered`
	tree := avltree.NewAVLTreeOrderedKey[int, string]()

	// Check tree is empty
	fmt.Println(tree.Empty()) // Prints true

	//Prints elements count
	fmt.Println(tree.Size()) // Prints 0

	// Insert a set of keys
	for i := 0; i <= 100; i += 10 {
		tree.Insert(i, strconv.Itoa(i))
	}
	// Now tree contains [0, 10, ... 100] keys with corresponded values

	// Check tree is empty
	fmt.Println(tree.Empty()) // Prints false

	//Prints elements count
	fmt.Println(tree.Size()) // Prints 11
}
```
</details>

## TODO
+ It will be interesting to provide some performance testing for this implementation and comparing with other not only Go.

## Links
