[![Go Reference](https://pkg.go.dev/badge/github.com/OlexiyKhokhlov/avltree.svg)](https://pkg.go.dev/github.com/OlexiyKhokhlov/avltree)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/OlexiyKhokhlov/avltree)](https://goreportcard.com/report/github.com/OlexiyKhokhlov/avltree)
# AVL Tree Go's module

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
+ Every tree node contains only two pointers. It helps to reduce memory usage but it blocks possibility to implement iterators. So iterators isn't present but you can use enumerator methods instead. I guess it is completely usefull.
+  Go doesn't allow ordinary developers to write generics like Go's map or slices. So this implentations uses interface{} type a lot. For you it means - you must cast interface{} to key or value every time when you want to reach your data. Also it reduces a performance since interally I have to cast interfaces in the key processing.
+ Go hasn't `const` qualifier. So there is no flex way to block possibility to change a key inside of the tree. Of course I know about copying. But isn't a good solution. First of all Go hasn't got unified way to copy any type of data. And secondly it provides a bad performance when your key is a big structure. So be carefull, avoid key changing! 
+ Regarding Go's interface{} again - you can insert a different types for values data. I don't recomend to do it but it will works correctly in any time. Just keep in the mind - correct casting to the value type back is completely on you! Also you can the tree like a set without values. In the such case you can use nil for every value.
+ Inserting and Erasing methods are non-recursive. So it helps to reduse stack memory usage. And also it is cool!

## Tutorial

### Creating a container
First of all you should implement a `Comparator`. This is a function that gets two keys via their interface{} and returns an int value that means:
+ 0 if both arguments are equivalent
+ -1 if the first argument is lesser then second
+ 1 if the first argument is greater then second
  
Example #1 Where key type is `int`:
```
func IntComparator func(a interface{}, b interface{}) int  {
	first := a.(int) //cast to key type here
	second := b.(int) //and here

	if first == second {
	    return 0
	}
	if first < second {
	    return -1
	}
	return 1
    }
```
Example #2 where key type is a custom struct:
```
type MyKey struct {
    field1 int
    key    uint64 //assume we want to use only this field like a key
    title  string
}

func MyKeyComparator func(a interface{}, b interface{}) int  {
	first := a.(MyKey) //cast to key type here
	second := b.(MyKey) //and here

	if first.key == second.key {
	    return 0
	}
	if first.key < second.key {
	    return -1
	}
	return 1
    }
```
And Example #3 Where key is a string:
```
func StringComparator func(a interface{}, b interface{}) int  {
	first := a.(string) //cast to key type here
	second := b.(string) //and here

	return first.Compare(second) //Use Compare method since it doing all needed
    }
```
And when you got a `Comparator` you can create a tree instance in a simple way:
```
package main

import (
	"github.com/OlexiyKhokhlov/avltree"
	"strings"
)

func StringComparator(a interface{}, b interface{}) int {
	first := a.(string)  //cast to key type here
	second := b.(string) //and here

	return strings.Compare(first, second) //Use Compare method since it doing all needed
}
func main() {
	StringTree := avltree.NewAVLTree(StringComparator)
}
}
```
Or create a tree where key has `int` type:
```
package main

import (
    "github.com/OlexiyKhokhlov/avltree"
)

func main() {
IntTree :=  avltree.NewAVLTree(func(a interface{}, b interface{}) int {
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
}
```

### Inserting and Erasing
Well, if you have a tree you want to insert and erase some data.
```
package main

import (
	"fmt"
	"github.com/OlexiyKhokhlov/avltree"
)

func main() {
	IntTree := avltree.NewAVLTree(func(a interface{}, b interface{}) int {
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

	//Insert a value here
	err := IntTree.Insert(10, nil)
	if err != nil {
		fmt.Println(err)
	}

	//Or insert a variable
	v := 20
	err = IntTree.Insert(v, nil)
	if err != nil {
		fmt.Println(err)
	}

	//The same key isn't allowed
	err = IntTree.Insert(20, nil)
	if err != nil {
		fmt.Println(err)
	}

	//Now IntTree contains 10 and 20
	//Lets try to remove 20
	err = IntTree.Erase(20)
	if err != nil {
		fmt.Println(err)
	}

	//Lets try to do the same again
	err = IntTree.Erase(20)
	if err != nil {
		fmt.Println(err)
	}

	//At the end clear all data in the tree
	IntTree.Clear()
}
```
### Element access, checking
```
package main

import (
	"fmt"
	"github.com/OlexiyKhokhlov/avltree"
)

func main() {
	IntTree := avltree.NewAVLTree(func(a interface{}, b interface{}) int {
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

	//Insert some data
	IntTree.Insert(1, "one")
	IntTree.Insert(2, "two")
	IntTree.Insert(3, "three")

	//Check if a tree contains some key
	if IntTree.Contains(1) {
		fmt.Println("IntTree truly contains 10")
	}

	//Find a value that is associated with a key
	value := IntTree.Find(2)
	if value != nil {
		str := value.(string)
		fmt.Println("Value for key 2 is: ", str)
	}
}
```
### Enumerating
TODO

## TODO
+ It will be interesting to provide some performance testing for this implementation and comparing with other not only Go.
+ Implemet a generator that generates a Tree sources by the given types.

## Links
