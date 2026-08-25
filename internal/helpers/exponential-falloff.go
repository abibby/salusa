package helpers

import (
	"math"
	"time"
)

func ExponentialFalloff(n int) time.Duration {
	return time.Duration(math.Min(math.Pow(2, float64(n)), 60)) * time.Second
}
