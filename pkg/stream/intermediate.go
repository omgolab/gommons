package gcstream

func (s *stream[T]) Filter(predicate PredicateFunc[T]) *stream[T] {
	newData := []T{}
	for _, v := range s.data {
		if predicate(v) {
			newData = append(newData, v)
		}
	}
	return &stream[T]{data: newData}
}

func (s *stream[T]) Map(mapper MapperFunc[T]) *stream[T] {
	for i, v := range s.data {
		s.data[i] = mapper(v)
	}
	return s
}
