package avltree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const MIN = -1000
const MAX = 1000

func fillTree(tree *AVLTree) {
	for i := MIN; i <= MAX; i++ {
		tree.Insert(i, nil)
	}
}

func TestCreate(t *testing.T) {
	require := require.New(t)

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

	require.Equal(empty_tree.Size(), uint(0))
	require.True(empty_tree.Empty())
}

func TestFindNode(t *testing.T) {
	require := require.New(t)

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

	const (
		START  = 0
		FINISH = 100
		STEP   = 5
	)
	for i := START; i <= FINISH; i += STEP {
		tree.Insert(i, nil)
	}

	require.Equal(START, tree.FirstNode(ASCENDING).Key)
	require.Equal(FINISH, tree.FirstNode(DESCENDING).Key)

	require.Nil(tree.FindPrevNode(START))
	require.Nil(tree.FindNextNode(FINISH))

	for i := START + 1; i <= FINISH; i += STEP {
		require.Equal(i-1, tree.FindPrevNode(i).Key)   //1
		require.Equal(i-1, tree.FindPrevNode(i+1).Key) //2
		require.Equal(i-1, tree.FindPrevNode(i+2).Key) //3
		require.Equal(i-1, tree.FindPrevNode(i+3).Key) //4
		require.Equal(i-1, tree.FindPrevNode(i+4).Key) //5
	}

	for i := FINISH - 1; i >= START; i -= STEP {
		require.Equal(i+1, tree.FindNextNode(i).Key)   //99
		require.Equal(i+1, tree.FindNextNode(i-1).Key) //98
		require.Equal(i+1, tree.FindNextNode(i-2).Key) //97
		require.Equal(i+1, tree.FindNextNode(i-3).Key) //96
		require.Equal(i+1, tree.FindNextNode(i-4).Key) //95
	}
}

func TestInsert(t *testing.T) {
	require := require.New(t)
	require.Less(MIN, MAX)

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

	var size uint = 0
	for i := MIN; i <= MAX; i++ {
		require.False(tree.Contains(i))
		require.Nil(tree.Insert(i, 0))
		require.Error(tree.Insert(i, 0))
		require.True(tree.Contains(i))
		size++
		require.Equal(tree.Size(), size)

		tree.checkHeight(func(hl int, hr int) {
			require.LessOrEqual(abs(hl-hr), 1)
		})
	}
}

func TestContains(t *testing.T) {
	require := require.New(t)

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

	fillTree(tree)

	require.False(tree.Contains(MIN - 1))
	require.False(tree.Contains(MAX + 1))
	require.False(tree.Contains(MIN * 2))
	require.False(tree.Contains(MAX * 2))

	for i := MIN; i <= MAX; i++ {
		require.True(tree.Contains(i))
	}
}

func TestErase(t *testing.T) {
	require := require.New(t)

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

	fillTree(tree)

	var size uint = tree.Size()
	for i := MIN; i <= MAX; i++ {
		require.True(tree.Contains(i))
		require.Nil(tree.Erase(i))
		require.False(tree.Contains(i))
		require.Error(tree.Erase(i))

		size--
		require.Equal(tree.Size(), size)

		tree.checkHeight(func(hl int, hr int) {
			require.LessOrEqual(abs(hl-hr), 1)
		})
	}
}

func TestEnumerate(t *testing.T) {
	require := require.New(t)

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

	fillTree(tree)

	i := MIN
	tree.Enumerate(ASCENDING, func(k interface{}, v interface{}) bool {
		require.Equal(k.(int), i)
		i++
		return true
	})

	i = MAX
	tree.Enumerate(DESCENDING, func(k interface{}, v interface{}) bool {
		require.Equal(k.(int), i)
		i--
		return true
	})

	i = MIN
	expected_interrupt := MIN + 10
	tree.Enumerate(ASCENDING, func(k interface{}, v interface{}) bool {
		require.Equal(k.(int), i)
		if k.(int) == expected_interrupt {
			return false
		}
		i++
		return true
	})
	require.Equal(i, expected_interrupt)

	i = MAX
	expected_interrupt = MAX - 10
	tree.Enumerate(DESCENDING, func(k interface{}, v interface{}) bool {
		require.Equal(k.(int), i)
		if k.(int) == expected_interrupt {
			return false
		}
		i--
		return true
	})
	require.Equal(i, expected_interrupt)
}
