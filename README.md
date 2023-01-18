# Iter: simple golang iterator using generics

## example1: simple iteration

```go
s := strings.Split("동해물과 백두산이 마르고 닳도록", " ")
fmt.Printf("%v\n", s)

it := iter.Of(s...)
for v, ok := it.Next(); ok; v, ok = it.Next() {
    fmt.Printf("%s\n", v)
}
```

output:

```txt
[동해물과 백두산이 마르고 닳도록]
동해물과
백두산이
마르고
닳도록
```

## example 2: filter and to slice

```go
s := []int{1, 2, 3, 4, 5, 6}
fmt.Printf("%v\n", s)

it := iter.Of(s...)
r := it.
    Filter(iter.Even[int]). // 2,4,6
    Slice()                 // [2, 4, 6]

fmt.Printf("%v\n", r)
```

output:

```txt
[1 2 3 4 5 6]
[2 4 6]
```

## example 3: filter, map and reduce

```go
s := []int{1, 2, 3, 4, 5, 6}
fmt.Printf("%v\n", s)

it := iter.Of(s...)
r := it.
    Filter(iter.Even[int]). // 2,4,6
    Map(iter.Multiply(2)).  // 4, 8, 12
    Reduce(iter.Add[int])   // 4 + 8 + 12

fmt.Printf("%v\n", r)
```

output:

```txt
[1 2 3 4 5 6]
24
```

## chan iteration

```go
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
```

output

```txt
[1 2 3 4 5 6]
[1 2 3 4 5 6]
```

## benchmark

```txt
goos: darwin
goarch: amd64
pkg: github.com/whitekid/iter
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkWordCount/single-16                 1 1079390286 ns/op  8811840 B/op      235 allocs/op
BenchmarkWordCount/goroutine-16              7  151540820 ns/op 19696067 B/op    32406 allocs/op
BenchmarkWordCount/iterator-16               7  153911633 ns/op 19717408 B/op    32477 allocs/op
```
