package g3

import (
	"math"
)

const (
	Pi      = float32(math.Pi)
	MathMax = math.MaxFloat32
)

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Tan(x float32) float32 {
	return float32(math.Tan(float64(x)))
}

func Deg2Rad(d float32) float32 {
	return d * math.Pi / 180.0
}

func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func MinMax(a, b float32) (min, max float32) {
	if a < b {
		return a, b
	}
	return b, a
}

func Clamp(a, min, max float32) float32 {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}

