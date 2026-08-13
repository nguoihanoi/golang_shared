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
