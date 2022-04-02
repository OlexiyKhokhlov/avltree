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

func fillTree(tree *AVLTree[int, interface{}]) {
	for i := MIN; i <= MAX; i++ {
		tree.Insert(i, nil)
	}
}

func TestCreate(t *testing.T) {
	require := require.New(t)

	emptyTree := NewAVLTreeOrderedKey[int, interface{}]()

	require.Equal(uint(0), emptyTree.Size())
	require.True(emptyTree.Empty())
}

func TestInsert(t *testing.T) {
	require := require.New(t)
	require.Less(MIN, MAX)

	tree := NewAVLTreeOrderedKey[int, interface{}]()

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

	tree := NewAVLTreeOrderedKey[int, interface{}]()

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

	tree := NewAVLTreeOrderedKey[int, interface{}]()

	const (
		START  = 0
		FINISH = 100
		STEP   = 5
	)
	for i := START; i <= FINISH; i += STEP {
		tree.Insert(i, nil)
	}

	key, _ := tree.First()
	require.Equal(START, *key)

	key, _ = tree.Last()
	require.Equal(FINISH, *key)

	require.Nil(tree.FindPrevElement(START))
	require.Nil(tree.FindNextElement(FINISH))

	for i := START + 1; i <= FINISH; i += STEP {
		k, _ := tree.FindPrevElement(i)
		require.Equal(i-1, *k) //1
		k, _ = tree.FindPrevElement(i + 1)
		require.Equal(i-1, *k) //2
		k, _ = tree.FindPrevElement(i + 2)
		require.Equal(i-1, *k) //3
		k, _ = tree.FindPrevElement(i + 3)
		require.Equal(i-1, *k) //4
		k, _ = tree.FindPrevElement(i + 4)
		require.Equal(i-1, *k) //5
	}

	for i := FINISH - 1; i >= START; i -= STEP {
		k, _ := tree.FindNextElement(i)
		require.Equal(i+1, *k) //99
		k, _ = tree.FindNextElement(i - 1)
		require.Equal(i+1, *k) //98
		k, _ = tree.FindNextElement(i - 2)
		require.Equal(i+1, *k) //97
		k, _ = tree.FindNextElement(i - 3)
		require.Equal(i+1, *k) //96
		k, _ = tree.FindNextElement(i - 4)
		require.Equal(i+1, *k) //95
	}
}

func TestErase(t *testing.T) {
	require := require.New(t)

	tree := NewAVLTreeOrderedKey[int, interface{}]()

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

	tree := NewAVLTreeOrderedKey[int, interface{}]()

	fillTree(tree)

	i := MIN
	tree.Enumerate(ASCENDING, func(k int, v interface{}) bool {
		require.Equal(i, k)
		i++
		return true
	})

	i = MAX
	tree.Enumerate(DESCENDING, func(k int, v interface{}) bool {
		require.Equal(i, k)
		i--
		return true
	})

	i = MIN
	expectedInterrupt := MIN + 10
	tree.Enumerate(ASCENDING, func(k int, v interface{}) bool {
		require.Equal(i, k)
		if k == expectedInterrupt {
			return false
		}
		i++
		return true
	})
	require.Equal(i, expectedInterrupt)

	i = MAX
	expectedInterrupt = MAX - 10
	tree.Enumerate(DESCENDING, func(k int, v interface{}) bool {
		require.Equal(i, k)
		if k == expectedInterrupt {
			return false
		}
		i--
		return true
	})
	require.Equal(expectedInterrupt, i)
}

func TestEnumerateDiapason(t *testing.T) {
	require := require.New(t)

	tree := NewAVLTreeOrderedKey[int, interface{}]()

	var (
		START  int = 0
		FINISH int = 100
		STEP   int = 5
	)

	for i := START; i <= FINISH; i += STEP {
		tree.Insert(i, nil)
	}

	//0..100
	i := START
	tree.EnumerateDiapason(nil, nil, ASCENDING, func(k int, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(nil, nil, DESCENDING, func(k int, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i - STEP
		return true
	})

	//0..100
	i = START
	tree.EnumerateDiapason(&START, &FINISH, ASCENDING, func(k int, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	tree.EnumerateDiapason(&START, &FINISH, DESCENDING, func(k int, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i - STEP
		return true
	})

	//0..100
	i = START
	left := START - 10
	right := FINISH + 10
	tree.EnumerateDiapason(&left, &right, ASCENDING, func(k int, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i + STEP
		return true
	})

	//100..0
	i = FINISH
	left = START - 10
	right = FINISH + 10
	tree.EnumerateDiapason(&left, &right, DESCENDING, func(k int, v interface{}) bool {
		require.True(i >= START && i <= FINISH)
		require.Equal(i, k)
		i = i - STEP
		return true
	})

	for i := START; i <= FINISH; i = i + STEP {
		j := i
		tree.EnumerateDiapason(&i, nil, ASCENDING, func(k int, v interface{}) bool {
			require.True(j >= i && j <= FINISH)
			require.Equal(j, k)
			j = j + STEP
			return true
		})

		l := FINISH
		tree.EnumerateDiapason(&l, nil, DESCENDING, func(k int, v interface{}) bool {
			require.True(l >= i && l <= FINISH)
			require.Equal(l, k)
			l = l - STEP
			return true
		})
	}

	for i := START; i <= FINISH; i = i + STEP {
		j := START
		tree.EnumerateDiapason(nil, &i, ASCENDING, func(k int, v interface{}) bool {
			require.True(j >= START && j <= i)
			require.Equal(j, k)
			j = j + STEP
			return true
		})

		l := i
		tree.EnumerateDiapason(nil, &i, DESCENDING, func(k int, v interface{}) bool {
			require.True(l >= START && l <= i)
			require.Equal(l, k)
			l = l - STEP
			return true
		})
	}
}

func TestMaxHeight(t *testing.T) {
	require := require.New(t)

	tree := NewAVLTreeOrderedKey[int, interface{}]()

	data := []int{33,
		20, 46,
		12, 28, 41, 51,
		07, 17, 25, 31, 38, 44, 49, 53,
		04, 10, 15, 19, 23, 27, 30, 32, 36, 40, 43, 45, 48, 50, 52,
		2, 6, 9, 11, 14, 16, 18, 22, 24, 26, 29, 35, 37, 39, 42, 47,
		1, 3, 5, 8, 13, 21, 34,
		0}
	for i := range data {
		tree.Insert(i, nil)
	}

	i := 0
	tree.Enumerate(ASCENDING, func(k int, v interface{}) bool {
		require.Equal(i, k)
		i++
		return true
	})
}
