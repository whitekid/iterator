package iter

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWordCount(t *testing.T) {
	data := loadData(t)
	want := wordCount(data)

	type args struct {
		countWords func(require.TestingT, io.Reader) int
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"iterator", args{countWithIterator}},
		{"single", args{countWithSingle}},
		{"goroutine", args{countWithGoroutine}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.countWords(t, bytes.NewReader(data))
			require.Equal(t, want, got)
		})
	}
}

func countWithSingle(t require.TestingT, r io.Reader) int {
	data, _ := io.ReadAll(r)
	return wordCount(data)
}

const chunkSize = 1000

func countWithGoroutine(t require.TestingT, r io.Reader) (count int) {
	countC := make(chan int)

	// map
	go func() {
		var wg sync.WaitGroup
		scanner := bufio.NewScanner(r)

		buf := make([]string, 0, chunkSize)
		for scanner.Scan() {
			buf = append(buf, scanner.Text())
			if len(buf) == chunkSize {
				wg.Add(1)
				go func(buf []string) {
					defer wg.Done()
					countC <- wordCount([]byte(strings.Join(buf, "\n")))
				}(buf)
				buf = make([]string, 0, chunkSize)
			}
		}

		if len(buf) > 0 {
			countC <- wordCount([]byte(strings.Join(buf, "\n")))
		}

		wg.Wait()
		defer close(countC)
	}()

	// reduce
	for c := range countC {
		count += c
	}

	return count
}

func countWithIterator(t require.TestingT, r io.Reader) int {
	it := Map(Chunk(splitLine(r), chunkSize), func(x []string) int {
		return wordCount([]byte(strings.Join(x, "\n")))
	})
	return Reduce(it, Add[int])
}

func BenchmarkWordCount(b *testing.B) {
	want := math.MinInt

	type args struct {
		countWords func(require.TestingT, io.Reader) int
	}

	benchmarks := [...]struct {
		name string
		args args
	}{
		{"single", args{countWithSingle}},
		{"goroutine", args{countWithGoroutine}},
		{"iterator", args{countWithIterator}},
	}
	for _, bb := range benchmarks {
		b.Run(bb.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := bytes.NewReader(loadData(b))
				b.StartTimer()

				got := bb.args.countWords(b, r)
				if want != math.MinInt {
					require.Equal(b, want, got)
				}
				want = got
			}
		})
	}
}

func wordCount(data []byte) (r int) {
	for i := 0; i < 100; i++ { // to make takes more time
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			r++
		}
	}

	return
}

func loadData(t require.TestingT) []byte {
	fixture := "fixtures/4300.txt"
	if _, err := os.Stat(fixture); os.IsNotExist(err) {
		resp, err := http.Get("https://raw.githubusercontent.com/bwhite/dv_hadoop_tests/master/python-streaming/word_count/input/4300.txt")
		require.NoError(t, err)
		defer resp.Body.Close()

		f, err := os.Create(fixture)
		require.NoError(t, err)

		io.Copy(f, resp.Body)
		f.Close()
	}
	data, err := os.ReadFile(fixture)
	require.NoError(t, err)
	return data
}
