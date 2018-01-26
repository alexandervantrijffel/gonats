package eventsourcing

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func random(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := min + int(math.Max(float64(r.Intn(max)-min-1), 0))%max
	fmt.Printf("waiting for %d seconds\n", result)
	return result
}
