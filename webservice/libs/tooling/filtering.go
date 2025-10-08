package tooling

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func HasAny[T any](arrayToBeFiltered []T, testFunc func(T) bool) bool {
	return len(Filter(arrayToBeFiltered, func(v T) bool { return testFunc(v) })) > 0
}
