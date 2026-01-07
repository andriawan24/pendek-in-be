package utils

func GetOrElse[T any](arg *T, defaultValue T) T {
	if arg != nil {
		return *arg
	} else {
		return defaultValue
	}
}
