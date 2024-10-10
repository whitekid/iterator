package iter

import (
	"iter"
	"maps"
)

type MapIterator[K comparable, V any] iter.Seq2[K, V]

func M[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				return
			}
		}
	}
}

func items[K comparable, V any](m map[K]V) iter.Seq2[K, V] { return maps.All(m) }

// func mapEach[K comparable, V any](m MapIterator[K, V], each func(K, V)) {
// 	fanOut(m.Items(), func(item Item[K, V]) { each(item.Key, item.Value) })
// }
