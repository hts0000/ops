package util

type Base interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string |
		~bool
}

// GetOrZero 返回指针对应的值，如果指针为nil，则返回指针对应类型的零值
func GetOrZero[T Base](ptr *T) T {
	if ptr != nil {
		return *ptr
	}
	var zero T
	return zero
}
