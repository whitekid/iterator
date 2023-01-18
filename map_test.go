package iter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

func TestMap(t *testing.T) {
	m := map[int]string{
		1: "first",
		2: "second",
		3: "third",
	}
	require.Equal(t,
		[]int{1, 2, 3},
		Sorted(M(m).Keys()).Slice())
	require.Equal(t,
		[]string{"first", "second", "third"},
		Sorted(M(m).Values()).Slice())
	require.Equal(t, []Item[int, string]{
		{Key: 3, Value: "third"},
		{Key: 2, Value: "second"},
		{Key: 1, Value: "first"},
	},
		SortedFunc(M(m).Items(),
			func(a, b Item[int, string]) bool { return Descending(a.Key, b.Key) }).
			Slice())
}

func TestMapKeys(t *testing.T) {
	type args struct {
		m map[int]string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{map[int]string{
			1: "first",
			2: "second",
			3: "third",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]int, 0, len(tt.args.m))
			for k := range tt.args.m {
				want = append(want, k)
			}
			slices.Sort(want)

			got := M(tt.args.m).Keys()
			require.Equal(t, want, Sorted(got).Slice())
		})
	}
}

func TestMapKeysChanClosed(t *testing.T) {
	type args struct {
		m map[int]string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{map[int]string{
			1: "first",
			2: "second",
			3: "third",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := M(tt.args.m).Keys().Next()
			require.True(t, ok)
		})
	}
}

func TestMapValues(t *testing.T) {
	type args struct {
		m map[int]string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{map[int]string{
			1: "first",
			2: "second",
			3: "third",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]string, 0, len(tt.args.m))
			for _, v := range tt.args.m {
				want = append(want, v)
			}
			slices.Sort(want)

			got := M(tt.args.m).Values()
			require.Equal(t, want, Sorted(got).Slice())
		})
	}
}

func TestMapItems(t *testing.T) {
	type args struct {
		m map[int]string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{`valid`, args{map[int]string{
			1: "first",
			2: "second",
			3: "third",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]Item[int, string], 0, len(tt.args.m))
			for k, v := range tt.args.m {
				want = append(want, Item[int, string]{k, v})
			}
			slices.SortFunc(want, func(a, b Item[int, string]) bool { return Asending(a.Key, b.Key) })

			got := M(tt.args.m).Items()
			require.Equal(t, want, SortedFunc(got, func(a, b Item[int, string]) bool { return Asending(a.Key, b.Key) }).Slice())
		})
	}
}
