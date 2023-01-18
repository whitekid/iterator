package iter

import (
	"math"
	"sync"
)

func S[S ~[]T, T any](s S) Iterator[T] {
	index := 0
	return &withNext[T]{
		next: func() (r T, ok bool) {
			if index >= len(s) {
				return r, false
			}

			index++
			return s[index-1], true
		},
	}
}
func Of[T any](s ...T) Iterator[T] { return S(s) }

func Reverse[T any](it Iterator[T]) Iterator[T] {
	index := math.MinInt
	var s []T
	o := sync.Once{}

	return &withNext[T]{
		next: func() (r T, ok bool) {
			if s != nil && index <= 0 {
				return r, false
			}

			o.Do(func() {
				s = slice(it)
				index = len(s)
			})

			index--
			return s[index], true
		},
	}
}

func Chunk[T any](it Iterator[T], size int) Iterator[[]T] {
	last := false

	return &withNext[[]T]{
		next: func() ([]T, bool) {
			chunk := make([]T, 0, size)
			i := 0
			for ; i < size; i++ {
				v, ok := it.Next()
				if !ok {
					break
				}
				chunk = append(chunk, v)
			}

			if last {
				return nil, false
			}

			last = len(chunk) != size
			return chunk, true
		},
	}
}
