package Math

import "github.com/go-gl/mathgl/mgl32"

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

func Mat3ToMat4(mat mgl32.Mat3) mgl32.Mat4 {
	return mgl32.Mat4{
		mat[0], mat[1], 0, mat[2],
		mat[3], mat[4], 0, mat[5],
		0, 0, 1, 0,
		mat[6], mat[7], 0, mat[8],
	}
}
