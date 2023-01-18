package iter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChan(t *testing.T) {
	ch := make(chan int)
	want := []int{1, 2, 3, 4, 5, 6}
	go func() {
		defer close(ch)
		for _, v := range want {
			ch <- v
		}
	}()

	it := C(ch)
	got := make([]int, 0)
	for v, ok := it.Next(); ok; v, ok = it.Next() {
		got = append(got, v)
	}

	require.Equal(t, want, got)
}
