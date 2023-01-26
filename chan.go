package iter

func C[T any](ch <-chan T) Iterator[T] {
	return &withNext[T]{
		next: func() (v T, ok bool) {
			v, ok = <-ch
			return v, ok
		},
	}
}

func fanOut[T any](it Iterator[T], fn func(T)) {
	q := newQueue[chan T]()

	go func() {
		defer q.Close()
		for v, ok := it.Next(); ok; v, ok = it.Next() {
			ch := q.Push(make(chan T))
			v := v
			go func() {
				ch <- v
				close(ch)
			}()
		}
	}()

	for ch, ok := q.Pop(); ok; ch, ok = q.Pop() {
		v := <-ch
		fn(v)
	}
}
