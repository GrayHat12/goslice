package outofplace

import (
	"context"

	"github.com/GrayHat12/goslice/commons"
	"github.com/GrayHat12/goslice/inplace"
)

func Filter[T any](slice []T, predicate func(*T, int, *[]T) bool) *[]T {
	indices := make([]int, 0, len(slice))

	for index, element := range slice {
		if !predicate(&element, index, &slice) {
			indices = append(indices, index)
		}
	}

	return inplace.RemoveManyElementsByIndices(&slice, indices)
}

func Map[T any, V any](slice []T, predicate func(*T, int, *[]T) *V) []V {
	outputmap := make([]V, len(slice))
	commons.ForEach(context.Background(), &slice, func(ctx context.Context, x *T, index int, array *[]T) {
		select {
		case <-ctx.Done():
			return
		default:
			outputmap[index] = *predicate(x, index, array)
		}
	})
	return outputmap
}
