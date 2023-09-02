package gcstream

// This stream package is a simple implementation of a stream in Go
// Similar Java stream: https://docs.oracle.com/en/java/javase/20/docs/api/java.base/java/util/stream/Stream.html
// TODO: create the intermediate methods lazy until a terminal method is called
// TODO: add parallel stream
// TODO: add more functions
// TODO: add tests
type stream[T any] struct {
	data []T
}

type PredicateFunc[T any] func(T) bool
type MapperFunc[T any] func(T) T
type ReducerFunc[T any] func(T, T) T

func New[T any](d []T) *stream[T] {
	return &stream[T]{data: d}
}
