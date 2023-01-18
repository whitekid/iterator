package iter

import (
	"runtime"
)

type Queue[T any] struct {
	items chan T
}

func newQueue[T any](size ...int) *Queue[T] {
	s := runtime.NumCPU()
	for _, e := range size {
		s = e
	}

	return &Queue[T]{
		items: make(chan T, s),
	}
}

func (q *Queue[T]) Close()     { close(q.items) }
func (q *Queue[T]) Push(v T) T { q.items <- v; return v }
func (q *Queue[T]) Pop() (T, bool) {
	v, ok := <-q.items
	return v, ok
}
