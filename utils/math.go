package utils

import "math"

// if v=6.36 return false if v=6.000 return true
func IsAnInt(v float64) bool {
	return math.Mod(v, 1) == 0
}
