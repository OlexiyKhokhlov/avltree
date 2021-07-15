package avltree

import (
	"errors"
	"fmt"
	"io"
)

// Function type that whould be defined for a key type in the tree
type Comparator func(a interface{}, b interface{}) int

type Enumerator func(key interface{}, value interface{}) bool

/// Internall stuff
func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

type node struct {
	Key     interface{}
	Value   interface{}
	Links   [2]*node
	Balance int
}

func getHeight(n *node) int {
	if n == nil {
		return 0
	}
	return max(getHeight(n.Links[0]), getHeight(n.Links[1])) + 1
}

func (n *node) getDirection(key interface{}, cmp Comparator) int {
	if cmp(key, n.Key) == -1 {
		return 0
	}
	return 1
}

func (n *node) isBalanced() bool {
	return n.Balance < 0
}

func lookup(n *node, key interface{}, cmp Comparator) *node {
	for n != nil && cmp(key, n.Key) != 0 {
		n = n.Links[n.getDirection(key, cmp)]
	}
	return n
}

func recursiveDump(n *node, w io.Writer) {
	if n != nil {
		io.WriteString(w, fmt.Sprintf("\"%v\"-> { ", n.Key))
		if n.Links[0] != nil {
			io.WriteString(w, fmt.Sprintf("\"%v\" ", n.Links[0].Key))
		}
		if n.Links[1] != nil {
			io.WriteString(w, fmt.Sprintf("\"%v\" ", n.Links[1].Key))
		}
		io.WriteString(w, "}\n")
		recursiveDump(n.Links[0], w)
		recursiveDump(n.Links[1], w)
	}
}

type heightChecker func(lh int, rh int)

func recursiveCheckHeight(n *node, checker heightChecker) {
	if n == nil {
		return
	}
	recursiveCheckHeight(n.Links[0], checker)
	recursiveCheckHeight(n.Links[1], checker)
	checker(getHeight(n.Links[0]), getHeight(n.Links[1]))
}

func rotate2(path_top **node, dir int) *node {
	node_B := *path_top
	node_D := node_B.Links[dir]
	node_C := node_D.Links[1-dir]
	node_E := node_D.Links[dir]

	*path_top = node_D
	node_D.Links[1-dir] = node_B
	node_B.Links[dir] = node_C

	return node_E
}

func rotate3(path_top **node, dir int) {
	node_B := *path_top
	node_F := node_B.Links[dir]
	node_D := node_F.Links[1-dir]
	/* note: C and E can be nil */
	node_C := node_D.Links[1-dir]
	node_E := node_D.Links[dir]
	*path_top = node_D
	node_D.Links[1-dir] = node_B
	node_D.Links[dir] = node_F
	node_B.Links[dir] = node_C
	node_F.Links[1-dir] = node_E
}

func avlRotate2(path_top **node, dir int) *node {
	(*path_top).Balance = -1
	result := rotate2(path_top, dir)
	(*path_top).Balance = -1
	return result
}

func avlRotate3(path_top **node, dir int, third int) *node {
	node_B := *path_top
	node_F := node_B.Links[dir]
	node_D := node_F.Links[1-dir]
	/* note: C and E can be nil */
	node_C := node_D.Links[1-dir]
	node_E := node_D.Links[dir]

	node_B.Balance = -1
	node_F.Balance = -1
	node_D.Balance = -1

	rotate3(path_top, dir)

	if third == -1 {
		return nil
	} else if third == dir {
		/* E holds the insertion so B is unbalanced */
		node_B.Balance = 1 - dir
		return node_E
	} else {
		/* C holds the insertion so F is unbalanced */
		node_F.Balance = dir
		return node_C
	}
}

func avlInsert(root **node, key interface{}, value interface{}, cmp Comparator) bool {
	//Stage 1. Find a position in the tree and link a new node
	// by the way find and remember a node where the tree starts to be unbalanced.
	node_ptr := root
	path_top := root
	n := *root
	for n != nil && cmp(key, n.Key) != 0 {
		if !n.isBalanced() {
			path_top = node_ptr
		}
		dir := n.getDirection(key, cmp)
		node_ptr = &(n.Links[dir])
		n = *node_ptr
	}
	if n != nil {
		return false //already has the key
	}
	new_node := &node{
		Key:     key,
		Value:   value,
		Balance: -1,
	}
	*node_ptr = new_node

	//Stage 2. Rebalance
	path := *path_top
	var first, second, third int
	if !path.isBalanced() {
		first = path.getDirection(key, cmp)
		if path.Balance != first {
			/* took the shorter path */
			path.Balance = -1
			path = path.Links[first]
		} else {
			second = path.Links[first].getDirection(key, cmp)
			if first == second {
				/* just a two-point rotate */
				path = avlRotate2(path_top, first)
			} else {
				/* fine details of the 3 point rotate depend on the third step.
				 * However there may not be a third step, if the third point of the
				 * rotation is the newly inserted point.  In that case we record
				 * the third step as NEITHER
				 */
				path = path.Links[first].Links[second]
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
		path.Balance = direction
		path = path.Links[direction]
	}
	return true
}

func avlErase(root **node, key interface{}, cmp Comparator) *node {
	//Stage 1. lookup for the node that contain a key
	n := *root
	nodep := root
	path_top := root
	var targetp **node
	var dir int

	for n != nil {
		dir = n.getDirection(key, cmp)
		if cmp(n.Key, key) == 0 {
			targetp = nodep
		} else if n.Links[dir] == nil {
			break
		} else if n.isBalanced() || (n.Balance == (1-dir) && n.Links[1-dir].isBalanced()) {
			path_top = nodep
		}
		nodep = &n.Links[dir]
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
	for true {
		tree := *treep
		bdir := tree.getDirection(key, cmp)
		if tree.Links[bdir] == nil {
			break
		} else if tree.isBalanced() {
			tree.Balance = 1 - bdir
		} else if tree.Balance == bdir {
			tree.Balance = -1
		} else {
			second := tree.Links[1-bdir].Balance
			if second == bdir {
				avlRotate3(treep, 1-bdir, tree.Links[1-bdir].Links[bdir].Balance)
			} else if second == -1 {
				avlRotate2(treep, 1-bdir)
				tree.Balance = 1 - bdir
				(*treep).Balance = bdir
			} else {
				avlRotate2(treep, 1-bdir)
			}
			if tree == targetn {
				targetp = &(*treep).Links[bdir]
			}
		}
		treep = &(tree.Links[bdir])
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
	*treep = tree.Links[1-dir]
	tree.Links[0] = targetn.Links[0]
	tree.Links[1] = targetn.Links[1]
	tree.Balance = targetn.Balance

	return targetn
}

/// Internall stuff END

type AVLTree struct {
	root    *node
	count   uint
	compare Comparator
}

func NewAVLTree(c Comparator) *AVLTree {
	return &AVLTree{
		compare: c,
	}
}

func (t *AVLTree) Size() uint {
	return t.count
}

func (t *AVLTree) Empty() bool {
	return t.count == 0
}

func (t *AVLTree) Contains(key interface{}) bool {
	return lookup(t.root, key, t.compare) != nil
}

func (t *AVLTree) Find(key interface{}) interface{} {
	n := lookup(t.root, key, t.compare)
	if n != nil {
		return n.Value
	}
	return nil
}

func (t *AVLTree) Insert(key interface{}, value interface{}) error {
	if avlInsert(&t.root, key, value, t.compare) {
		t.count++
		return nil
	}
	return errors.New("AVLTree: already contains key")
}

func (t *AVLTree) Erase(key interface{}) error {
	if nil != avlErase(&t.root, key, t.compare) {
		t.count--
		return nil
	}
	return errors.New("AVLTree: hasn't got key")
}

func (t *AVLTree) Clear() {
	t.root = nil
	t.count = 0
}

func recursiveEnumAsk(n *node, f Enumerator) {
	if n.Links[0] != nil {
		recursiveEnumAsk(n.Links[0], f)
	}
	f(n.Key, n.Value)
	if n.Links[1] != nil {
		recursiveEnumAsk(n.Links[1], f)
	}
}

func recursiveEnumDesc(n *node, f Enumerator) {
	if n.Links[1] != nil {
		recursiveEnumDesc(n.Links[1], f)
	}
	f(n.Key, n.Value)
	if n.Links[0] != nil {
		recursiveEnumDesc(n.Links[0], f)
	}
}

func (t *AVLTree) EnumerateAsc(f Enumerator) {
	if t.root != nil {
		recursiveEnumAsk(t.root, f)
	}
}

func (t *AVLTree) EnumerateDesc(f Enumerator) {
	if t.root != nil {
		recursiveEnumDesc(t.root, f)
	}
}

func (t *AVLTree) BSTDump(w io.Writer) {
	io.WriteString(w, "digraph BST {\n")
	recursiveDump(t.root, w)
	io.WriteString(w, "}\n")
}

func (t *AVLTree) CheckHeight(checker heightChecker) {
	recursiveCheckHeight(t.root, checker)
}
