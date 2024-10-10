package iter

func fanOut[T any](it Iterator[T], fn func(T)) {
	q := newQueue[chan T]()

	go func() {
		defer q.Close()
		for v := range it {
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
