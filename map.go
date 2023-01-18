package iter

import (
	"golang.org/x/exp/maps"
)

type MapIterator[K comparable, V any] interface {
	Keys() Iterator[K]
	Values() Iterator[V]
	Items() Iterator[Item[K, V]]
}

type Item[K comparable, V any] struct {
	Key   K
	Value V
}

type mapIter[K comparable, V any] struct {
	orig map[K]V
}

func M[K comparable, V any](m map[K]V) MapIterator[K, V] {
	return &mapIter[K, V]{
		orig: m,
	}
}

func (m *mapIter[K, V]) Keys() Iterator[K]            { return keys(m.orig) }
func keys[K comparable, V any](m map[K]V) Iterator[K] { return S(maps.Keys(m)) }

func (m *mapIter[K, V]) Values() Iterator[V]            { return values(m.orig) }
func values[K comparable, V any](m map[K]V) Iterator[V] { return S(maps.Values(m)) }

func (m *mapIter[K, V]) Items() Iterator[Item[K, V]] { return items(m.orig) }
func items[K comparable, V any](m map[K]V) Iterator[Item[K, V]] {
	q := newQueue[Item[K, V]]()
	go func() {
		defer q.Close()
		for k, v := range m {
			q.Push(Item[K, V]{k, v})
		}
	}()
	return &withNext[Item[K, V]]{
		next: func() (Item[K, V], bool) {
			item, ok := q.Pop()

			return item, ok
		},
	}
}
