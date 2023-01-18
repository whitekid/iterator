package main

import (
	"fmt"
	"strings"

	"github.com/whitekid/iter"
)

func main() {
	run("example1: simple iteration", example1)
	run("example2: filter and to slice", example2)
	run("example3: filter, map and reduce", example3)
	run("example chan", exampleChan)
}

func run(s string, fn func()) {
	fmt.Printf("===== %s ====\n", s)
	fn()
}

func example1() {
	s := strings.Split("동해물과 백두산이 마르고 닳도록", " ")
	fmt.Printf("%v\n", s)

	it := iter.Of(s...)
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		fmt.Printf("%s\n", v)
	}
}

func example2() {
	s := []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%v\n", s)

	it := iter.Of(s...)
	r := it.
		Filter(iter.Even[int]). // 2,4,6
		Slice()                 // [2, 4, 6]

	fmt.Printf("%v\n", r)
}

func example3() {
	s := []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%v\n", s)

	it := iter.Of(s...)
	r := it.
		Filter(iter.Even[int]). // 2,4,6
		Map(iter.Multiply(2)).  // 4, 8, 12
		Reduce(iter.Add[int])   // 4 + 8 + 12

	fmt.Printf("%v\n", r)
}

func exampleChan() {
	s := []int{1, 2, 3, 4, 5, 6}
	fmt.Printf("%v\n", s)

	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, e := range s {
			ch <- e
		}
	}()

	it := iter.C(ch)
	r := make([]int, 0)
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		r = append(r, v)
	}

	fmt.Printf("%v\n", r)
}
