package util

import "golang.org/x/exp/constraints"

// Min will find the minimum of two numbers
func Min[T constraints.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// Max will find the maximum of two numbers
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
