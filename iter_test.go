package iter

import (
	"bufio"
	"bytes"
	"io"
	"iter"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// check for goroutine
	// goleak.VerifyTestMain(m)
	os.Exit(m.Run())
}

func TestMapper(t *testing.T) {
	type args struct {
		s      []int
		mapper func(int) int
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{[]int{1, 2, 3, 4, 5}, Multiply(2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]int, len(tt.args.s))
			for i, x := range tt.args.s {
				want[i] = tt.args.mapper(x)
			}

			it := Map(slices.Values(tt.args.s), tt.args.mapper)
			got := Collect(it)
			require.Equal(t, got, tt.args.s)
		})
	}
}

func testMapSignle(t require.TestingT, r io.Reader) {
	scanner := bufio.NewScanner(r)
	items := make([]int, 0)
	for scanner.Scan() {
		items = append(items, wordCount([]byte(scanner.Text())))
	}
	_ = items
}

func testMap(t require.TestingT, r io.Reader) {
	it := Map(splitLine(r), func(s string) int { return wordCount([]byte(s)) })
	_ = Collect(it)
}

func splitLine(r io.Reader) iter.Seq[string] {
	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				break
			}
		}
	}
}

func BenchmarkMapper(b *testing.B) {
	type args struct {
		mapper func(require.TestingT, io.Reader)
	}
	bencharmks := [...]struct {
		name string
		args args
	}{
		{"single", args{testMapSignle}},
		{"iterator", args{testMap}},
	}
	for _, bb := range bencharmks {
		b.Run(bb.name, func(b *testing.B) {
			b.StopTimer()
			data := loadData(b)
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				bb.args.mapper(b, bytes.NewReader(data))
			}
		})
	}
}

func TestReduce(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	sum := it.Reduce(Add[int])
	require.Equal(t, 15, sum)
}

func TestFilter(t *testing.T) {
	type args struct {
		s      []int
		filter func(int) bool
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{[]int{1, 2, 3, 4, 5}, Even[int]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]int, 0, len(tt.args.s))
			for _, x := range tt.args.s {
				if tt.args.filter(x) {
					want = append(want, x)
				}
			}

			got := Of(tt.args.s...).Filter(tt.args.filter).Collect()
			require.Equal(t, want, got)
		})
	}
}

func TestSlice_(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := Collect(slices.Values(s))
	require.Equal(t, s, got)
}

func TestSliceWithFilter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := Collect(Filter(slices.Values(s), Even[int]))
	require.Equal(t, []int{2, 4}, got)
}

func TestSort(t *testing.T) {
	s := []int{3, 2, 4, 1, 2, 5}
	require.Equal(t, Sorted(slices.Values(s)), SortedFunc(slices.Values(s), Asending[int]))
	require.Equal(t, Reverse(Sorted(slices.Values(s))), SortedFunc(slices.Values(s), Descending[int]))
}

func TestConcat(t *testing.T) {
	it := Concat(Of(1, 2, 3).Seq(), Of(4, 5, 6).Seq(), Of(7, 8, 9).Seq())
	got := Collect(it)
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
}

func TestTakeWhile(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	got := it.TakeWhile(func(x int) bool { return x < 4 }).Collect()
	require.Equal(t, []int{1, 2, 3}, got)
}

func TestDropWhile(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	got := it.DropWhile(func(x int) bool { return x < 3 }).Collect()
	require.Equal(t, []int{3, 4, 5}, got)
}

func TestSkip(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	got := it.Skip(3).Collect()
	require.Equal(t, []int{4, 5}, got)
}

func TestEach(t *testing.T) {
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	got := []int{}
	Of(want...).Each(func(x int) { got = append(got, x) })
	require.Equal(t, want, got)
}

func TestEachIdx(t *testing.T) {
	want := []string{"one", "two", "three", "four", "five"}

	got := []int{}
	Of(want...).EachIdx(func(i int, x string) { got = append(got, i) })
	require.Equal(t, []int{0, 1, 2, 3, 4}, got)
}
