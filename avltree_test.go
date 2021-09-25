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

	require.Equal(uint(0), empty_tree.Size())
	require.True(empty_tree.Empty())
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
		require.Equal(size, tree.Size())

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

func TestFindExt(t *testing.T) {
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

	key, _ := tree.First()
	require.Equal(START, key)

	key, _ = tree.Last()
	require.Equal(FINISH, key)

	require.Nil(tree.FindPrevElement(START))
	require.Nil(tree.FindNextElement(FINISH))

	for i := START + 1; i <= FINISH; i += STEP {
		k, _ := tree.FindPrevElement(i)
		require.Equal(i-1, k) //1
		k, _ = tree.FindPrevElement(i + 1)
		require.Equal(i-1, k) //2
		k, _ = tree.FindPrevElement(i + 2)
		require.Equal(i-1, k) //3
		k, _ = tree.FindPrevElement(i + 3)
		require.Equal(i-1, k) //4
		k, _ = tree.FindPrevElement(i + 4)
		require.Equal(i-1, k) //5
	}

	for i := FINISH - 1; i >= START; i -= STEP {
		k, _ := tree.FindNextElement(i)
		require.Equal(i+1, k) //99
		k, _ = tree.FindNextElement(i - 1)
		require.Equal(i+1, k) //98
		k, _ = tree.FindNextElement(i - 2)
		require.Equal(i+1, k) //97
		k, _ = tree.FindNextElement(i - 3)
		require.Equal(i+1, k) //96
		k, _ = tree.FindNextElement(i - 4)
		require.Equal(i+1, k) //95
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
		require.Equal(size, tree.Size())

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
		require.Equal(i, k.(int))
		i++
		return true
	})

	i = MAX
	tree.Enumerate(DESCENDING, func(k interface{}, v interface{}) bool {
		require.Equal(i, k.(int))
		i--
		return true
	})

	i = MIN
	expected_interrupt := MIN + 10
	tree.Enumerate(ASCENDING, func(k interface{}, v interface{}) bool {
		require.Equal(i, k.(int))
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
		require.Equal(i, k.(int))
		if k.(int) == expected_interrupt {
			return false
		}
		i--
		return true
	})
	require.Equal(expected_interrupt, i)
}

func TestEnumerateDiapason(t *testing.T) {
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

	//0..100
	i := START
	tree.EnumerateDiapason(nil, nil, ASCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k.(int))
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(nil, nil, DESCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k.(int))
		i = i - STEP
		return true
	})

	//0..100
	i = START
	tree.EnumerateDiapason(START, FINISH, ASCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k.(int))
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(START, FINISH, DESCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k.(int))
		i = i - STEP
		return true
	})

	//0..100
	i = START
	tree.EnumerateDiapason(START-10, FINISH+10, ASCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k.(int))
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(START-10, FINISH+10, DESCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k.(int))
		i = i - STEP
		return true
	})

	//0..20
	i = START
	right := 20
	tree.EnumerateDiapason(nil, right, ASCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= right)
		require.Equal(i, k.(int))
		i = i + STEP
		return true
	})

	//20..0
	i = START
	right = 21
	tree.EnumerateDiapason(nil, right, ASCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= right)
		require.Equal(i, k.(int))
		i = i + STEP
		return true
	})

	//65..100
	i = 65
	left := 63
	tree.EnumerateDiapason(left, nil, ASCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= 60 && i <= FINISH)
		require.Equal(i, k.(int))
		i = i + STEP
		return true
	})

	//100..65
	i = 65
	right = 67
	tree.EnumerateDiapason(nil, right, DESCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= right)
		require.Equal(i, k.(int))
		i = i - STEP
		return true
	})

	//65..0
	i = 65
	right = 67
	tree.EnumerateDiapason(nil, right, DESCENDING, func(k interface{}, v interface{}) bool {
		require.True(i >= START && i <= right)
		require.Equal(i, k.(int))
		i = i - STEP
		return true
	})
}
