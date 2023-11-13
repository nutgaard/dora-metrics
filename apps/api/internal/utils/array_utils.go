package utils

func Map[T any, O any](elements []T, mapper func(element T, i int) O) []O {
	result := make([]O, 0, len(elements))
	for i, element := range elements {
		result = append(result, mapper(element, i))
	}
	return result
}
func Filter[T any](elements []T, predicate func(element T, i int) bool) []T {
	result := make([]T, 0, len(elements))
	for i, element := range elements {
		add := predicate(element, i)
		if add {
			result = append(result, element)
		}
	}
	return result
}
func Reduce[A any, T any](elements []T, initialValue A, mapper func(acc A, value T, i int) A) A {
	var acc = initialValue
	for i, element := range elements {
		acc = mapper(acc, element, i)
	}
	return acc
}
