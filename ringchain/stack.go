package ringchain

type stack[T comparable] struct {
	elements []T
	registry map[T]struct{}
}

func newStack[T comparable]() *stack[T] {
	return &stack[T]{
		elements: make([]T, 0),
		registry: make(map[T]struct{}),
	}
}

func (s *stack[T]) push(t T) {
	s.elements = append(s.elements, t)
	s.registry[t] = struct{}{}
}

func (s *stack[T]) pop() (T, bool) {
	element, ok := s.top()
	if !ok {
		return element, false
	}

	s.elements = s.elements[:len(s.elements)-1]
	delete(s.registry, element)

	return element, true
}

func (s *stack[T]) top() (T, bool) {
	if s.isEmpty() {
		var defaultValue T
		return defaultValue, false
	}

	return s.elements[len(s.elements)-1], true
}

func (s *stack[T]) isEmpty() bool {
	return len(s.elements) == 0
}
