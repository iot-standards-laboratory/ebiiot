package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func GetSleepTime() int {
	// return rand.ExpFloat64()*10 + 95
	return (1 + rand.Intn(10)) * 50
}
