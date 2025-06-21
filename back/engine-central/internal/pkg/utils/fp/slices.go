package fp

func Map[A, B any](xs []A, f func(A, int) B) []B {
	ys := make([]B, len(xs))
	for i, x := range xs {
		ys[i] = f(x, i)
	}
	return ys
}

func Filter[A any](f func(A) bool, xs []A) []A {

	ys := make([]A, 0)
	for _, x := range xs {
		if f(x) {
			ys = append(ys, x)
		}
	}
	return ys
}
