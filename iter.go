package iter

import (
	"iter"
	"slices"
	"strconv"

	"golang.org/x/exp/constraints"
)

func Of[T any](s ...T) Iterator[T] {
	return Iterator[T](slices.Values(s))
}

type Iterator[T any] iter.Seq[T]

func (it Iterator[T]) Seq() iter.Seq[T] { return iter.Seq[T](it) }

func (it Iterator[T]) Map(fn func(T) T) Iterator[T] {
	return Iterator[T](Map[T](iter.Seq[T](it), fn))
}

func (it Iterator[T]) Filter(fn func(T) bool) Iterator[T] {
	return Iterator[T](Filter(iter.Seq[T](it), fn))
}

func (it Iterator[T]) TakeWhile(fn func(T) bool) Iterator[T] {
	return Iterator[T](TakeWhile(iter.Seq[T](it), fn))
}

func (it Iterator[T]) DropWhile(fn func(T) bool) Iterator[T] {
	return Iterator[T](DropWhile(iter.Seq[T](it), fn))
}

func (it Iterator[T]) Skip(n int) Iterator[T] {
	return Iterator[T](Skip[T](iter.Seq[T](it), n))
}
func (it Iterator[T]) Reduce(fn func(T, T) T) T { return Reduce[T](iter.Seq[T](it), fn) }
func (it Iterator[T]) Collect() []T             { return Collect[T](iter.Seq[T](it)) }
func (it Iterator[T]) Each(fn func(T))          { Each[T](iter.Seq[T](it), fn) }
func (it Iterator[T]) EachIdx(fn func(int, T))  { eachIdx[T](it, fn) }

// Map map mapper using goroutine
func Map[T1, T2 any](it iter.Seq[T1], mapper func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v := range it {
			if !yield(mapper(v)) {
				break
			}
		}
	}
}

// Sample map functions
func StrToInt(s string) (v int)               { v, _ = strconv.Atoi(s); return }
func Multiply[T Number](factor T) func(x T) T { return func(x T) T { return x * factor } }

type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

func Filter[T any](s iter.Seq[T], filterer func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := range s {
			if filterer(e) {
				if !yield(e) {
					break
				}
			}
		}
	}
}

// Sample filter functions
func Even[T constraints.Integer](x T) bool { return x%2 == 0 }
func Odd[T constraints.Integer](x T) bool  { return x%2 == 1 }

func Reduce[T any](it iter.Seq[T], reducer func(T, T) T) T {
	var v T

	for e := range it {
		v = e
		break
	}

	for e := range it {
		v = reducer(v, e)
	}

	return v
}

func TakeWhile[T any](it iter.Seq[T], take func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !take(v) {
				break
			}

			if !yield(v) {
				break
			}
		}
	}
}

func DropWhile[T any](it iter.Seq[T], drop func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := range it {
			if drop(e) {
				continue
			}

			if !yield(e) {
				return
			}
		}

		for e := range it {
			if !yield(e) {
				break
			}
		}
	}
}

func Skip[T any](it iter.Seq[T], n int) iter.Seq[T] {
	for i := 0; i < n; i++ {
		for _ = range it {
			break
		}
	}

	return it
}

// Sample reducer functions
func Add[T Number | string](a, b T) T { return a + b }

func Sum[T Number | string](it iter.Seq[T]) T { return Reduce(it, Add[T]) }

func Each[T any](it iter.Seq[T], each func(T)) {
	for v := range it {
		each(v)
	}
}

func eachIdx[T any](it Iterator[T], each func(int, T)) {
	i := 0
	fanOut(it, func(x T) {
		each(i, x)
		i++
	})
}

func Sorted[T constraints.Ordered](it iter.Seq[T]) iter.Seq[T] {
	return Of(slices.Sorted(it)...).Seq()
}

type Less[T any] func(T, T) int

func Asending[T constraints.Integer | constraints.Float](a, b T) int   { return int(a - b) }
func Descending[T constraints.Integer | constraints.Float](a, b T) int { return -Asending(a, b) }

func SortedFunc[T any](it iter.Seq[T], less Less[T]) []T {
	return slices.SortedFunc(it, less)
}

func Concat[T any](it ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, i := range it {
			for e := range i {
				if !yield(e) {
					return
				}
			}
		}
	}
}
