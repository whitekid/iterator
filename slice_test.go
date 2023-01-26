package iter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	{
		require.Equalf(t, []int{1, 2, 3, 4, 5, 6}, Of(s...).Slice(), "slice failed()")
		require.Equalf(t, []int{2, 4, 6}, Of(s...).Filter(Even[int]).Slice(), "slice/filter failed")
		require.Equalf(t, []int{1, 3, 5}, Of(s...).Filter(Odd[int]).Slice(), "slice/failter failed")
		require.Equalf(t, []int{4, 8, 12}, Of(s...).Filter(Even[int]).Map(Multiply(2)).Slice(), "slice/failter failed")
	}

	{
		it := Of(s...)
		for i, v := range s {
			value, ok := it.Next()
			require.Truef(t, ok, "index: %v", i)
			require.Equal(t, v, value, "index: %v", i)
		}

		_, ok := it.Next()
		require.False(t, ok)
	}

	{
		got := Of(s...).Reduce(Add[int])
		require.Equal(t, 21, got)
	}

	{
		want := []int{}
		for _, v := range s {
			want = append(want, v*2)
		}

		it := Of(s...)
		it1 := Map(it, func(x int) int { return x * 2 })

		got := []int{}
		for v, ok := it1.Next(); ok; v, ok = it1.Next() {
			got = append(got, v)
		}

		require.Equal(t, want, got)
	}

	{
		it := Of(s...)
		it1 := Map(it, func(x int) int { return x * 2 })
		sum := reduce(it1, Add[int])

		want := 0
		for _, v := range s {
			want += v * 2
		}
		require.Equalf(t, want, sum, "s=%v", s)

		require.Equal(t, 21, Sum(Of(s...)))
	}
}

func TestReverse(t *testing.T) {
	type args struct {
		s []int
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{[]int{1, 2, 3, 4, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]int, len(tt.args.s))
			copy(want, tt.args.s)
			slices.SortFunc(want, Descending[int])

			it := Reverse(S(tt.args.s))
			got := slice(it)

			require.Equal(t, want, got)
		})
	}
}

func TestChunk(t *testing.T) {
	type args struct {
		s    []int
		size int
	}

	tests := [...]struct {
		name string
		arg  args
		want [][]int
	}{
		{"valid", args{[]int{1, 2, 3, 4, 5}, 2}, [][]int{{1, 2}, {3, 4}, {5}}},
		{`valid`, args{[]int{1, 2, 3, 4}, 2}, [][]int{{1, 2}, {3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{
				it := S(tt.arg.s)
				chunk := Chunk(it, tt.arg.size)
				for i := 0; i < len(tt.want); i++ {
					v, ok := chunk.Next()
					require.True(t, ok)
					require.Equal(t, tt.want[i], v)
				}
			}

			{
				it := S(tt.arg.s)
				got := Chunk(it, tt.arg.size).Slice()
				require.Equal(t, tt.want, got)
			}
		})
	}
}
