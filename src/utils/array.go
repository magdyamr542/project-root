package utils

func Filter(array []string, predicate func(string) bool) (ret []string) {
	for _, s := range array {
		if predicate(s) {
			ret = append(ret, s)
		}
	}
	return
}

func Contains(array []string, element string) bool {
	for _, value := range array {
		if value == element {
			return true
		}
	}
	return false
}

func Map(array []string, transformer func(string) string) []string {
	var result []string
	for _, s := range array {
		result = append(result, transformer(s))
	}
	return result
}
