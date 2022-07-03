package utilities

import (
	"math/rand"
	"time"
)

// SetRandomSeed as per golang documentation,
// it is enough to call this function once
// in the process to set the random seed
func SetRandomSeed() {
	rand.Seed(time.Now().UnixNano())
}
