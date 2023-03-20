package log

func MustGet[T any](value T, err error) func(l Logger) T {
	return func(l Logger) T {
		l.Must(err)
		return value
	}
}
