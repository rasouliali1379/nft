package lambda

func Map[T1, T2 any](in []T1, f func(T1) T2) []T2 {
	res := make([]T2, len(in))
	for i, item := range in {
		res[i] = f(item)
	}
	return res
}

func Filter[T any](in []T, f func(T) bool) []T {
	res := make([]T, 0, len(in))
	for _, item := range in {
		if f(item) {
			res = append(res, item)
		}
	}
	return res
}

func Reduce[T1, T2 any](in []T1, init T2, reducer func(T2, T1) T2) T2 {
	for _, item := range in {
		init = reducer(init, item)
	}
	return init
}

func Any[T any](in []T, f func(T) bool) bool {
	res := false
	for i := 0; i < len(in); i++ {
		res = res || f(in[i])
	}
	return res
}

func All[T any](in []T, f func(T) bool) bool {
	res := true
	for i := 0; i < len(in); i++ {
		res = res && f(in[i])
	}
	return res
}