package avltree

import (
	"testing"
)

func BenchmarkAVLInsert(b *testing.B) {
	tree := NewAVLTreeOrderedKey[int, interface{}]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1_000; j++ {
			tree.Insert(j, nil)
		}
	}
}

func BenchmarkMAPInsert(b *testing.B) {
	m := make(map[int]interface{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1_000; j++ {
			m[j] = nil
		}
	}
}
