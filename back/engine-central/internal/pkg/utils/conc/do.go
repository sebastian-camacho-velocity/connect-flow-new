package conc

// Do executes all the given functions in parallel and returns a MultiError
func Do(funcs ...func() error) error {
	return Each(len(funcs), funcs, func(in func() error) error {
		return in()
	})
}
