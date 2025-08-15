package helper

func IntPtrToFloat64(ptr *int) float64 {
	if ptr == nil {
		return 0
	}
	return float64(*ptr)
}
