package util

import "errors"

func Flatten2DArray[T interface{}](arr [][]T) []T {
	var flat []T

	for _, row := range arr {
		flat = append(flat, row...)
	}

	return flat
}

func Unflatten2DArray[T interface{}](arr []T, rowLength int) ([][]T, error) {
	if rowLength <= 0 {
		return nil, errors.New("row length must be greater than zero")
	}

	if len(arr)%rowLength != 0 {
		return nil, errors.New("flat array length is not a multiple of row length")
	}

	var result [][]T
	for i := 0; i < len(arr); i += rowLength {
		result = append(result, arr[i:i+rowLength])
	}

	return result, nil
}
