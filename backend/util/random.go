package util

import (
	rand "math/rand/v2"
)

func RandRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func RandIntSlice(minLength, maxLength, min, max int) []int {
	ret := make([]int, RandRange(minLength, maxLength))

	for i := range ret {
		ret[i] = RandRange(min, max)
	}

	return ret
}
