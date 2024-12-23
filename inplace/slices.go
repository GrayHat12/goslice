package inplace

func RemoveManyElementsByIndices[T any](slice *[]T, indices []int) *[]T {
	indicesMap := make(map[int]int)

	for _, index := range indices {
		indicesMap[index] = index
	}

	lastIndex := len(*slice) - 1
	backIndex := lastIndex

	for _, index := range indices {
		if index < 0 || index > lastIndex {
			continue
		}

		mappedIndex := indicesMap[index]

		if mappedIndex == -1 {
			continue
		}

		if mappedIndex != backIndex {
			(*slice)[mappedIndex] = (*slice)[backIndex]

			indicesMap[backIndex] = indicesMap[mappedIndex]
		}

		indicesMap[index] = -1

		backIndex--
	}

	newaddress := (*slice)[:backIndex+1]
	return &newaddress
}

func Filter[T any](slice *[]T, predicate func(*T, int, *[]T) bool) *[]T {
	indices := make([]int, 0, len(*slice))

	for index, element := range *slice {
		if !predicate(&element, index, slice) {
			indices = append(indices, index)
		}
	}

	return RemoveManyElementsByIndices(slice, indices)
}
