package iter

import (
	"iter"
	"slices"
)


func Collect[T any](s iter.Seq[T]) (r []T) {
	return slices.Collect(s)
}

func Reverse[T any](it iter.Seq[T]) iter.Seq[T] {
	s := Collect(it)

	return func(yield func(T) bool) {
		for i := len(s) - 1; i > 0; i-- {
			if !yield(s[i]) {
				return
			}
		}
	}
}

func Chunk[T any](it iter.Seq[T], size int) iter.Seq[[]T] {
	return slices.Chunk(Collect(it), size)
}
