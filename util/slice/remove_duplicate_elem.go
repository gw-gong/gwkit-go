package slice

func RemoveDuplicateElems[T comparable](slice []T) []T {
	elemSet := make(map[T]struct{})
	result := make([]T, 0)
	for _, elem := range slice {
		if _, ok := elemSet[elem]; !ok {
			elemSet[elem] = struct{}{}
			result = append(result, elem)
		}
	}
	return result
}
