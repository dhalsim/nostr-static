package utils

// Ternary is a generic function that returns a if cond is true, otherwise returns b
func Ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
