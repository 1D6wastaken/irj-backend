package utils

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func KeepError(_ any, err error) error {
	return err
}

func PtrTo[T any](t T) *T {
	return &t
}

func DerefOr[T any](v *T, def T) T {
	if v == nil {
		return def
	}

	return *v
}
