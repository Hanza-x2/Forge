package Math

func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func Norm(a, b, t float32) float32 {
	return (t - a) / (b - a)
}

func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
