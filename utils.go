package biligo

import (
	"context"
	"sync"
	"unsafe"
)

func toString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func await(tasks ...func()) error {
	return awaitCtx(context.Background(), tasks...)
}

func awaitCtx(ctx context.Context, tasks ...func()) error {
	wg := &sync.WaitGroup{}
	done := make(chan struct{})

	for _, f := range tasks {
		wg.Go(f)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type set[T comparable] map[T]struct{}

func (set *set[T]) add(v ...T) {
	for _, val := range v {
		(*set)[val] = struct{}{}
	}
}

func (set *set[T]) delete(v ...T) {
	for _, val := range v {
		delete(*set, val)
	}
}

// get 不保证顺序
func (set *set[T]) get() (s []T) {
	s = make([]T, len(*set))
	i := 0
	for val := range *set {
		s[i] = val
		i++
	}
	return
}

func (set *set[T]) ok(v T) bool {
	_, ok := (*set)[v]
	return ok
}

// clean 能保证顺序
func (set *set[T]) clean(s []T) []T {
	write := 0
	for read := range s {
		if !set.ok(s[read]) {
			(*set)[s[read]] = struct{}{}
			s[write] = s[read]
			write++
		}
	}
	return s[:write]
}
