package utilities

type arrayClass struct {
	//radius float64 // Private field
}

func Array() *arrayClass {
	return &arrayClass{}
}

func (a *arrayClass) RemoveDuplicates(elements []string, isEmpty bool) []string {
	encountered := make(map[string]struct{}, len(elements))
	result := make([]string, 0, len(elements))
	if isEmpty == false {
		for _, val := range elements {
			if _, exists := encountered[val]; !exists {
				encountered[val] = struct{}{}
				if val != "" {
					result = append(result, val)
				}
			}
		}
		return result
	}
	for _, val := range elements {
		if _, exists := encountered[val]; !exists {
			encountered[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}

func SymmetricDifference[T comparable](a, b []T) []T {
	seen := make(map[T]int)
	for _, v := range a {
		seen[v] |= 1
	}
	for _, v := range b {
		seen[v] |= 2
	}
	var result []T
	for v, mask := range seen {
		if mask == 1 || mask == 2 {
			result = append(result, v)
		}
	}
	return result
}

func Difference[T comparable](a, b []T) []T {
	bMap := make(map[T]struct{}, len(b))
	for _, v := range b {
		bMap[v] = struct{}{}
	}

	var result []T
	for _, v := range a {
		if _, found := bMap[v]; !found {
			result = append(result, v)
		}
	}

	return result
}
