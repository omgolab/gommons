package gcstream

func (s *stream[T]) AnyMatch(predicate PredicateFunc[T]) bool {
	for _, v := range s.data {
		if predicate(v) {
			return true
		}
	}
	return false
}

func (s *stream[T]) AllMatch(predicate PredicateFunc[T]) bool {
	for _, v := range s.data {
		if !predicate(v) {
			return false
		}
	}
	return true
}

func (s *stream[T]) Reduce(initial T, reducer ReducerFunc[T]) T {
	acc := initial
	for _, v := range s.data {
		acc = reducer(acc, v)
	}
	return acc
}

func (s *stream[T]) Collect() []T {
	return s.data
}
