package avltree

import "testing"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const MIN = -1000
const MAX = 1000

func TestCreateAVLTree(t *testing.T) {
	empty_tree := NewAVLTree(func(a interface{}, b interface{}) int {
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

	if size := empty_tree.Size(); size != 0 {
		t.Error("Expected Size 0, got ", size)
	}

	if empty := empty_tree.Empty(); !empty {
		t.Error("Expected Empty true, got ", empty)
	}
}

func eraseAVLTree(tree *AVLTree, t *testing.T) {
	var size uint = tree.Size()
	for i := MIN; i <= MAX; i++ {
		if !tree.Contains(i) {
			t.Error("Tree hasn't got ", i)
		}
		if tree.Erase(i) != nil {
			t.Error("Can't erase ", i)
			break
		}
		if tree.Contains(i) {
			t.Error("Key wasn't removed ", i)
			break
		}
		if tree.Erase(i) == nil {
			t.Error("Tree hasn't got a key but it has been removed ", i)
			break
		}
		size--
		if tree.Size() != size {
			t.Error(
				"expected size ", size,
				"got ", tree.Size(),
			)
			break
		}

		tree.CheckHeight(func(hl int, hr int) {
			if abs(hl-hr) > 1 {
				t.Error("Tree isn't balanced after inserting of ", i)
			}
		})
	}
}

func containsAVLTRee(tree *AVLTree, t *testing.T) {
	if tree.Contains(MIN - 1) {
		t.Error("Tree unexpectedly contains ", MIN-1)
	}
	if tree.Contains(MAX + 1) {
		t.Error("Tree unexpectedly contains ", MAX+1)
	}
	if tree.Contains(MAX * 2) {
		t.Error("Tree unexpectedly contains ", MAX*2)
	}
	if tree.Contains(MIN * 2) {
		t.Error("Tree unexpectedly contains ", MIN*2)
	}
	for i := MIN; i <= MAX; i++ {
		if !tree.Contains(i) {
			t.Error("Tree hasn't got ", i)
		}
	}

}

func insertAVLTree(tree *AVLTree, t *testing.T) {
	var size uint = 0
	for i := MIN; i <= MAX; i++ {
		if tree.Contains(i) {
			t.Error("Tree unexpectedly contains ", i)
		}
		if tree.Insert(i, 0) != nil {
			t.Error("Can't insert ", i)
			break
		}
		if tree.Insert(i, 0) == nil {
			t.Error("Duplicate insertion must be not allowed for ", i)
			break
		}
		if !tree.Contains(i) {
			t.Error("Tree hasn't got a key that has been inserted right now. ", i)
			break
		}
		size++
		if tree.Size() != size {
			t.Error(
				"expected size ", size,
				"got ", tree.Size(),
			)
			break
		}

		tree.CheckHeight(func(hl int, hr int) {
			if abs(hl-hr) > 1 {
				t.Error("Tree isn't balanced after inserting of ", i)
			}
		})
	}
}

func TestInsertErase(t *testing.T) {
	if MIN >= MAX {
		t.Error("WRONG MIN - MAX diapason")
	}

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
	insertAVLTree(tree, t)
	containsAVLTRee(tree, t)
	eraseAVLTree(tree, t)
}
