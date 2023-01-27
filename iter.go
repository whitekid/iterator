package iter

import (
	"strconv"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// NOTE Next()와 Value()가 thread safe하지 않음..
type Iterator[T any] interface {
	Next() (T, bool)

	Filter(func(T) bool) Iterator[T]
	Map(func(T) T) Iterator[T]
	// TODO Map(func(T, T) T) Iterator[T2] go1.19 does not supports yet
	TakeWhile(func(T) bool) Iterator[T]
	DropWhile(func(T) bool) Iterator[T]
	Skip(n int) Iterator[T]

	Reduce(func(T, T) T) T
	Slice() []T
	Each(func(T))
	EachIdx(func(int, T))
}

type withNext[T any] struct {
	next func() (T, bool)
}

func (it *withNext[T]) Next() (T, bool)                       { return it.next() }
func (it *withNext[T]) Map(fn func(T) T) Iterator[T]          { return Map[T](it, fn) }
func (it *withNext[T]) Filter(fn func(T) bool) Iterator[T]    { return filter[T](it, fn) }
func (it *withNext[T]) TakeWhile(fn func(T) bool) Iterator[T] { return takeWhile[T](it, fn) }
func (it *withNext[T]) DropWhile(fn func(T) bool) Iterator[T] { return dropWhile[T](it, fn) }
func (it *withNext[T]) Skip(n int) Iterator[T]                { return skip[T](it, n) }
func (it *withNext[T]) Reduce(fn func(T, T) T) T              { return reduce[T](it, fn) }
func (it *withNext[T]) Slice() (r []T)                        { return slice[T](it) }
func (it *withNext[T]) Each(fn func(T))                       { each[T](it, fn) }
func (it *withNext[T]) EachIdx(fn func(int, T))               { eachIdx[T](it, fn) }

// Map map mapper using goroutine
func Map[T1, T2 any](it Iterator[T1], mapper func(T1) T2) Iterator[T2] {
	q := newQueue[chan T2]()
	go func() {
		defer q.Close()
		for v, ok := it.Next(); ok; v, ok = it.Next() {
			ch := q.Push(make(chan T2))
			v := v
			go func() {
				ch <- mapper(v)
				close(ch)
			}()
		}
	}()

	return &withNext[T2]{
		next: func() (r T2, ok bool) {
			ch, ok := q.Pop()
			if ok {
				return <-ch, ok
			}
			return r, ok
		},
	}
}

// Sample map functions
func StrToInt(s string) (v int)               { v, _ = strconv.Atoi(s); return }
func Multiply[T Number](factor T) func(x T) T { return func(x T) T { return x * factor } }

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

func filter[T any](it Iterator[T], filterer func(T) bool) Iterator[T] {
	return &withNext[T]{
		next: func() (r T, ok bool) {
			for v, ok := it.Next(); ok; v, ok = it.Next() {
				if ok && filterer(v) {
					return v, ok
				}
			}
			return r, false
		},
	}
}

// Sample filter functions
func Even[T constraints.Integer](x T) bool { return x%2 == 0 }
func Odd[T constraints.Integer](x T) bool  { return x%2 == 1 }

func reduce[T any](it Iterator[T], reducer func(T, T) T) T {
	value, ok := it.Next()
	if !ok {
		return value
	}

	fanOut(it, func(v T) {
		value = reducer(value, v)
	})

	return value
}

func takeWhile[T any](it Iterator[T], take func(T) bool) Iterator[T] {
	q := newQueue[T]()
	go func() {
		defer q.Close()
		for v, ok := it.Next(); ok; v, ok = it.Next() {
			if !take(v) {
				break
			}
			q.Push(v)
		}
	}()

	return &withNext[T]{
		next: func() (T, bool) {
			v, ok := q.Pop()
			return v, ok
		},
	}
}

func dropWhile[T any](it Iterator[T], drop func(T) bool) Iterator[T] {
	q := newQueue[T]()
	go func() {
		defer q.Close()

		for v, ok := it.Next(); ok; v, ok = it.Next() {
			if !drop(v) {
				q.Push(v)
				break
			}
		}

		for v, ok := it.Next(); ok; v, ok = it.Next() {
			q.Push(v)
		}
	}()

	return &withNext[T]{
		next: func() (T, bool) {
			v, ok := q.Pop()
			return v, ok
		},
	}
}

func skip[T any](it Iterator[T], n int) Iterator[T] {
	for i := 0; i < n; i++ {
		it.Next()
	}

	return it
}

// Sample reducer functions
func Add[T Number | string](a, b T) T { return a + b }

func Sum[T Number | string](it Iterator[T]) T { return reduce(it, Add[T]) }

func slice[T any](it Iterator[T]) (r []T) {
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		r = append(r, v)
	}
	return r
}

func each[T any](it Iterator[T], each func(T)) {
	fanOut(it, func(x T) { each(x) })
}

func eachIdx[T any](it Iterator[T], each func(int, T)) {
	i := 0
	fanOut(it, func(x T) {
		each(i, x)
		i++
	})
}

func Sorted[T constraints.Ordered](it Iterator[T]) Iterator[T] {
	s := it.Slice()
	slices.Sort(s)
	return S(s)
}

type Less[T any] func(T, T) bool

func Asending[T constraints.Ordered](a, b T) bool   { return a < b }
func Descending[T constraints.Ordered](a, b T) bool { return !Asending(a, b) }

func SortedFunc[T any](it Iterator[T], less Less[T]) Iterator[T] {
	s := it.Slice()
	slices.SortFunc(s, less)
	return S(s)
}

func Concat[T any](it ...Iterator[T]) Iterator[T] {
	q := newQueue[T]()
	go func() {
		defer q.Close()
		for _, i := range it {
			for v, ok := i.Next(); ok; v, ok = i.Next() {
				q.Push(v)
			}
		}
	}()

	return &withNext[T]{
		next: func() (T, bool) {
			v, ok := q.Pop()
			return v, ok
		},
	}
}
