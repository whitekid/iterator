package iter

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// check for goroutine
	// goleak.VerifyTestMain(m)
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

			{
				it := Map(S(tt.args.s), tt.args.mapper)
				got := []int{}
				for i := 0; i < len(tt.args.s); i++ {
					v, ok := it.Next()
					require.True(t, ok)
					got = append(got, v)
				}

				_, ok := it.Next()
				require.False(t, ok, "next() want true but got false, want=%v, got=%v", want, got)
				require.Equalf(t, want, got, "want %v, got=%v", want, got)
			}

			{
				it := Map(S(tt.args.s), tt.args.mapper)
				got := slice(it)
				require.Equal(t, want, got)
			}
			{
				want := 0
				for _, x := range tt.args.s {
					want += tt.args.mapper(x)
				}

				it := Map(S(tt.args.s), tt.args.mapper)
				got := reduce(it, Add[int])
				require.Equal(t, want, got)
			}
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
	_ = slice(it)
}

func splitLine(r io.Reader) Iterator[string] {
	scanner := bufio.NewScanner(r)
	return &withNext[string]{
		next: func() (string, bool) {
			ok := scanner.Scan()
			text := scanner.Text()
			return text, ok
		},
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
	sum := Of(1, 2, 3, 4, 5).Reduce(Add[int])
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

			it := S(tt.args.s)
			it = filter(it, tt.args.filter)
			got := slice(it)
			require.Equal(t, want, got)
		})
	}
}

func TestSlice_(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := slice(S(s))
	require.Equal(t, s, got)
}

func TestSliceWithFilter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := slice(filter(S(s), Even[int]))
	require.Equal(t, []int{2, 4}, got)
}

func TestSort(t *testing.T) {
	s := []int{3, 2, 4, 1, 2, 5}
	require.Equal(t,
		Sorted(S(s)).Slice(), SortedFunc(S(s),
			Asending[int]).Slice())
	require.Equal(t,
		Reverse(Sorted(S(s))).Slice(), SortedFunc(S(s),
			Descending[int]).Slice())
}

func TestConcat(t *testing.T) {
	it := Concat(Of(1, 2, 3), Of(4, 5, 6), Of(7, 8, 9))
	got := it.Slice()
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
}

func TestTakeWhile(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	got := it.TakeWhile(func(x int) bool { return x < 4 }).Slice()
	require.Equal(t, []int{1, 2, 3}, got)
}

func TestDropWhile(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	got := it.DropWhile(func(x int) bool { return x < 3 }).Slice()
	require.Equal(t, []int{3, 4, 5}, got)
}

func TestSkip(t *testing.T) {
	it := Of(1, 2, 3, 4, 5)
	got := it.Skip(3).Slice()
	require.Equal(t, []int{4, 5}, got)
}
