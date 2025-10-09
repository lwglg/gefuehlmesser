package tooling

import "time"

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

func InTimeSpanEx(start, end, check time.Time, includeStart, includeEnd bool) bool {
	_start := start
	_end := end
	_check := check

	if end.Before(start) {
		_end = end.Add(24 * time.Hour)

		if check.Before(start) {
			_check = check.Add(24 * time.Hour)
		}
	}

	if includeStart {
		_start = _start.Add(-1 * time.Nanosecond)
	}

	if includeEnd {
		_end = _end.Add(1 * time.Nanosecond)
	}

	return _check.After(_start) && _check.Before(_end)
}
