// AVLTree are associative containers that store elements formed by a combination of a key value and a mapped value, following a specific order.
// In a AVLTree, the key values are generally used to sort and uniquely identify the elements, while the mapped values store the content associated to this key.
// The types of key and mapped value may differ.
// Internally, the elements in a AVLTree are always sorted by its key following a specific strict weak ordering criterion
// indicated by its internal comparison object (of type Comparator).
// AVLTree containers are generally slower than go map container to access individual elements by their key,
// but they allow the direct iteration on subsets based on their order.
package avltree

import (
	"errors"
	"fmt"
	"io"
	"math/bits"
	"unsafe"
)

// Function type that whould be defined for a key type in the tree
type Comparator func(a interface{}, b interface{}) int

// Function type for AVLTree traversal
type Enumerator func(key interface{}, value interface{}) bool

/// Internall stuff
func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

type Node struct {
	Key     interface{}
	Value   interface{}
	links   [2]*Node
	balance int
}

func getHeight(n *Node) int {
	if n == nil {
		return 0
	}
	return max(getHeight(n.links[0]), getHeight(n.links[1])) + 1
}

func (n *Node) getDirection(key interface{}, cmp Comparator) int {
	if cmp(key, n.Key) == -1 {
		return 0
	}
	return 1
}

func (n *Node) avlIsBalanced() bool {
	return n.balance < 0
}

func recursiveDump(n *Node, w io.Writer) {
	if n != nil {
		io.WriteString(w, fmt.Sprintf("\"%v\"-> { ", n.Key))
		if n.links[0] != nil {
			io.WriteString(w, fmt.Sprintf("\"%v\" ", n.links[0].Key))
		}
		if n.links[1] != nil {
			io.WriteString(w, fmt.Sprintf("\"%v\" ", n.links[1].Key))
		}
		io.WriteString(w, "}\n")
		recursiveDump(n.links[0], w)
		recursiveDump(n.links[1], w)
	}
}

type heightChecker func(lh int, rh int)

func recursiveCheckHeight(n *Node, checker heightChecker) {
	if n == nil {
		return
	}
	recursiveCheckHeight(n.links[0], checker)
	recursiveCheckHeight(n.links[1], checker)
	checker(getHeight(n.links[0]), getHeight(n.links[1]))
}

func (t *AVLTree) checkHeight(checker heightChecker) {
	recursiveCheckHeight(t.root, checker)
}

func rotate2(path_top **Node, dir int) *Node {
	node_B := *path_top
	node_D := node_B.links[dir]
	node_C := node_D.links[1-dir]
	node_E := node_D.links[dir]

	*path_top = node_D
	node_D.links[1-dir] = node_B
	node_B.links[dir] = node_C

	return node_E
}

func rotate3(path_top **Node, dir int) {
	node_B := *path_top
	node_F := node_B.links[dir]
	node_D := node_F.links[1-dir]
	/* note: C and E can be nil */
	node_C := node_D.links[1-dir]
	node_E := node_D.links[dir]
	*path_top = node_D
	node_D.links[1-dir] = node_B
	node_D.links[dir] = node_F
	node_B.links[dir] = node_C
	node_F.links[1-dir] = node_E
}

func avlRotate2(path_top **Node, dir int) *Node {
	(*path_top).balance = -1
	result := rotate2(path_top, dir)
	(*path_top).balance = -1
	return result
}

func avlRotate3(path_top **Node, dir int, third int) *Node {
	node_B := *path_top
	node_F := node_B.links[dir]
	node_D := node_F.links[1-dir]
	/* note: C and E can be nil */
	node_C := node_D.links[1-dir]
	node_E := node_D.links[dir]

	node_B.balance = -1
	node_F.balance = -1
	node_D.balance = -1

	rotate3(path_top, dir)

	if third == -1 {
		return nil
	} else if third == dir {
		/* E holds the insertion so B is unbalanced */
		node_B.balance = 1 - dir
		return node_E
	} else {
		/* C holds the insertion so F is unbalanced */
		node_F.balance = dir
		return node_C
	}
}

func avlInsert(root **Node, key interface{}, value interface{}, cmp Comparator) bool {
	//Stage 1. Find a position in the tree and link a new node
	// by the way find and remember a node where the tree starts to be unbalanced.
	node_ptr := root
	path_top := root
	n := *root
	for n != nil && cmp(key, n.Key) != 0 {
		if !n.avlIsBalanced() {
			path_top = node_ptr
		}
		dir := n.getDirection(key, cmp)
		node_ptr = &(n.links[dir])
		n = *node_ptr
	}
	if n != nil {
		return false //already has the key
	}
	new_node := &Node{
		Key:     key,
		Value:   value,
		balance: -1,
	}
	*node_ptr = new_node

	//Stage 2. Rebalance
	path := *path_top
	var first, second, third int
	if !path.avlIsBalanced() {
		first = path.getDirection(key, cmp)
		if path.balance != first {
			/* took the shorter path */
			path.balance = -1
			path = path.links[first]
		} else {
			second = path.links[first].getDirection(key, cmp)
			if first == second {
				/* just a two-point rotate */
				path = avlRotate2(path_top, first)
			} else {
				/* fine details of the 3 point rotate depend on the third step.
				 * However there may not be a third step, if the third point of the
				 * rotation is the newly inserted point.  In that case we record
				 * the third step as NEITHER
				 */
				path = path.links[first].links[second]
				if cmp(key, path.Key) == 0 {
					third = -1
				} else {
					third = path.getDirection(key, cmp)
				}
				path = avlRotate3(path_top, first, third)
			}
		}
	}

	//Stage 3. Update balance info in the each node
	for path != nil && cmp(key, path.Key) != 0 {
		direction := path.getDirection(key, cmp)
		path.balance = direction
		path = path.links[direction]
	}
	return true
}

func avlErase(root **Node, key interface{}, cmp Comparator) *Node {
	//Stage 1. lookup for the node that contain a key
	n := *root
	nodep := root
	path_top := root
	var targetp **Node
	var dir int

	for n != nil {
		dir = n.getDirection(key, cmp)
		if cmp(n.Key, key) == 0 {
			targetp = nodep
		} else if n.links[dir] == nil {
			break
		} else if n.avlIsBalanced() || (n.balance == (1-dir) && n.links[1-dir].avlIsBalanced()) {
			path_top = nodep
		}
		nodep = &n.links[dir]
		n = *nodep
	}
	if targetp == nil {
		return nil //key not found nothing to remove
	}

	/*
	 * Stage 2.
	 * adjust balance, but don't lose 'targetp'.
	 * each node from treep down towards target, but
	 * excluding the last, will have a subtree grow
	 * and need rebalancing
	 */
	treep := path_top
	targetn := *targetp
	for {
		tree := *treep
		bdir := tree.getDirection(key, cmp)
		if tree.links[bdir] == nil {
			break
		} else if tree.avlIsBalanced() {
			tree.balance = 1 - bdir
		} else if tree.balance == bdir {
			tree.balance = -1
		} else {
			second := tree.links[1-bdir].balance
			if second == bdir {
				avlRotate3(treep, 1-bdir, tree.links[1-bdir].links[bdir].balance)
			} else if second == -1 {
				avlRotate2(treep, 1-bdir)
				tree.balance = 1 - bdir
				(*treep).balance = bdir
			} else {
				avlRotate2(treep, 1-bdir)
			}
			if tree == targetn {
				targetp = &(*treep).links[bdir]
			}
		}
		treep = &(tree.links[bdir])
	}

	/*
	 * Stage 3.
	 * We have re-balanced everything, it remains only to
	 * swap the end of the path (*treep) with the deleted item
	 * (*targetp)
	 */
	tree := *treep
	targetn = *targetp
	*targetp = tree
	*treep = tree.links[1-dir]
	tree.links[0] = targetn.links[0]
	tree.links[1] = targetn.links[1]
	tree.balance = targetn.balance

	return targetn
}

func (t *AVLTree) findEdgeNodeImpl(key interface{}, dir int) *Node {
	var n, candidate *Node = t.root, nil
	for n != nil {
		cmp_res := t.compare(key, n.Key)
		if cmp_res == (2*dir - 1) {
			n = n.links[dir]
			continue
		}
		if cmp_res == (1 - 2*dir) {
			candidate = n
			n = n.links[1-dir]
			continue
		}
		// cmp_res == 0
		if n.links[dir] == nil {
			return candidate
		}
		return edgeNodeImpl(n.links[dir], 1-dir)
	}

	return candidate
}

func edgeNodeImpl(n *Node, dir int) *Node {
	for n.links[dir] != nil {
		n = n.links[dir]
	}
	return n
}

func (t *AVLTree) lookupNode(key interface{}) *Node {
	n := t.root
	for n != nil {
		cmp := t.compare(key, n.Key)
		if cmp == 0 {
			return n
		} else if cmp == -1 {
			n = n.links[0]
		} else {
			n = n.links[1]
		}
	}
	return n
}

/// Internall stuff END

// AVLTree is a sorted associative container that contains key-value pairs with unique keys.
// Keys are sorted by using the comparison function `Comparator`.
// Search, removal, and insertion operations have logarithmic complexity.
type AVLTree struct {
	root    *Node
	count   uint
	compare Comparator
}

// Creates a new AVLTree instance with the given Comparator
func NewAVLTree(c Comparator) *AVLTree {
	return &AVLTree{
		compare: c,
	}
}

// Returns the number of elements
func (t *AVLTree) Size() uint {
	return t.count
}

// Checks whether the container is empty
func (t *AVLTree) Empty() bool {
	return t.count == 0
}

// Checks if the container contains element with specific key
func (t *AVLTree) Contains(key interface{}) bool {
	return t.lookupNode(key) != nil
}

// Finds element with specific key
// Returns an interface{} for associated with the key value.
// When key isn't present returns nil
func (t *AVLTree) Find(key interface{}) interface{} {
	n := t.lookupNode(key)
	if n != nil {
		return n.Value
	}
	return nil
}

// Returns key and value that is nearest to the given key and lesser then given key.
// Can return (nil, nil) when no such node in the tree
func (t *AVLTree) FindPrevElement(key interface{}) (interface{}, interface{}) {
	node := t.findEdgeNodeImpl(key, 0)
	if node != nil {
		return node.Key, node.Value
	}
	return nil, nil
}

// Returns key and value with the key that is nearest to the given key and greater then given key.
// Can return (nil, nil) when no such node in the tree
func (t *AVLTree) FindNextElement(key interface{}) (interface{}, interface{}) {
	node := t.findEdgeNodeImpl(key, 1)
	if node != nil {
		return node.Key, node.Value
	}
	return nil, nil
}

// Inserts an element with the given key and value.
// Value can be nil
// It the given key is already present returns an error
func (t *AVLTree) Insert(key interface{}, value interface{}) error {
	if avlInsert(&t.root, key, value, t.compare) {
		t.count++
		return nil
	}
	return errors.New("AVLTree: already contains key")
}

// This a type of enumeration for Enumerate, EnumerateLowerBound, EnumerateUpperBound methods
// There are two acceptable values -  ASCENDING and DESCENDING. All other values provides a runtime error.
// Unfortunately Go doesn't provide any possibiltiy to check wrong values for that in the compile time. So be carefull here!
type EnumerationOrder int

const (
	ASCENDING  = 0
	DESCENDING = 1
)

// Returns key, value interfaces for the first tree node
// Returns (nil, nil) when a tree is empty
func (t *AVLTree) First() (interface{}, interface{}) {
	node := edgeNodeImpl(t.root, ASCENDING)
	if node == nil {
		return nil, nil
	}
	return node.Key, node.Value
}

// Returns key, value interfaces for the last tree node
// Returns (nil, nil) when a tree is empty
func (t *AVLTree) Last() (interface{}, interface{}) {
	node := edgeNodeImpl(t.root, DESCENDING)
	if node == nil {
		return nil, nil
	}
	return node.Key, node.Value
}

// Removes a element by the given key
func (t *AVLTree) Erase(key interface{}) error {
	if nil != avlErase(&t.root, key, t.compare) {
		t.count--
		return nil
	}
	return errors.New("AVLTree: key not found")
}

// clears the contents
func (t *AVLTree) Clear() {
	t.root = nil
	t.count = 0
}

// Calls 'Enumerator' for every Tree's element.
// Enumeration order can be one from ASCENDING or DESCENDING
// Enumerator should return `false` for stop enumerating or `true` for continue
func (t *AVLTree) Enumerate(order EnumerationOrder, f Enumerator) {
	n := t.root
	if n == nil {
		return
	}

	stack := make([]*Node, bits.Len(t.count))
	stack_ptr := 0
loop:
	for {
		switch uintptr(unsafe.Pointer(n)) & 0x01 {
		case 0: //Going down as deep as possible
			for ; n.links[order] != nil; n = n.links[order] {
				stack[stack_ptr] = (*Node)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) | 0x01))
				stack_ptr++
			}
			fallthrough
		case 1: //Going first up
			n = (*Node)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) & ^uintptr(0x01)))
			// Visit node
			if !f(n.Key, n.Value) {
				return
			}
			// Going down via second link or return up
			if next := n.links[1-order]; next != nil {
				n = next
			} else if stack_ptr != 0 {
				stack_ptr--
				n = stack[stack_ptr]
			} else {
				break loop
			}
		}
	}
}

// Works like Enumerate but has two additional args - left and right
// These are left and right borders for enumeration.
// Enumeration includes left and right borders.
// Note: left must be always lesser than right. Otherwise returns error
// Note: left and right should be nil. In means the lesser/greater key in the tree is a border.
//       So call EnumerateDiapason where both borders are nil is equivalent to call Enumerate.
// Note: If you want to enumerate whole tree call Enumerate since it`s faster!
func (t *AVLTree) EnumerateDiapason(left, right interface{}, order EnumerationOrder, f Enumerator) error {
	if t.count == 0 {
		return nil
	}

	if left != nil && right != nil && t.compare(left, right) > 0 {
		return errors.New("AVLTree: left must be less rigth")
	}

	//find common sub-tree
	n := t.root
	for {
		if left != nil && t.compare(n.Key, left) < 0 {
			n = n.links[1]
			continue
		}
		if right != nil && t.compare(n.Key, right) > 0 {
			n = n.links[0]
			continue
		}
		break
	}

	fences := [2]interface{}{left, right}
	stack := make([]*Node, bits.Len(t.count))
	stack_ptr := 0
loop:
	for {
		switch uintptr(unsafe.Pointer(n)) & 0x01 {
		case 0: //Going down
			for {
				if fences[order] != nil && (1-2*int(order))*t.compare(n.Key, fences[order]) < 0 {
					// Try go down via second link
					if next := n.links[1-order]; next != nil && (fences[1-order] == nil || (fences[1-order] != nil && (1-2*int(order))*t.compare(next.Key, fences[1-order]) <= 0)) {
						n = next
						continue
					} else if stack_ptr != 0 {
						//Or return up
						stack_ptr--
						n = stack[stack_ptr]
						continue loop
					} else {
						break loop
					}
				} else {
					if n.links[order] == nil {
						break
					}
					stack[stack_ptr] = (*Node)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) | 0x01))
					stack_ptr++
					n = n.links[order]
				}
			}
			fallthrough
		case 1: //Going first up
			n = (*Node)(unsafe.Pointer(uintptr(unsafe.Pointer(n)) & ^uintptr(0x01)))
			// Visit node
			if !f(n.Key, n.Value) {
				return nil
			}
			// Going down via second link
			if next := n.links[1-order]; next != nil {
				if fences[1-order] != nil && (1-2*int(order))*t.compare(next.Key, fences[1-order]) >= 0 {
					for ; next != nil; next = next.links[order] {
						if (1-2*int(order))*t.compare(next.Key, fences[1-order]) <= 0 {
							n = next
							continue loop
						}
					}
				} else {
					n = next
					continue loop
				}
			}
			if stack_ptr != 0 {
				stack_ptr--
				n = stack[stack_ptr]
			} else {
				break loop
			}
		}
	}
	return nil
}

// Writes BST Tree in graphviz digraph textual format
// See here https://graphviz.org/ for the details
func (t *AVLTree) BSTDump(w io.Writer) {
	io.WriteString(w, "digraph BST {\n")
	recursiveDump(t.root, w)
	io.WriteString(w, "}\n")
}
