package utils

func Filter(array []string, predicate func(string) bool) (ret []string) {
	for _, s := range array {
		if predicate(s) {
			ret = append(ret, s)
		}
	}
	return
}
