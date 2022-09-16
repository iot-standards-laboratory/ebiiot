package utils

import "sync"

var count int = 0
var mutex sync.Mutex

func Cummulate() int {
	mutex.Lock()
	defer mutex.Unlock()
	count++
	return count
}

var wg sync.WaitGroup

func Done() {
	wg.Done()
}

func Wait(count int) {
	wg.Add(count)
	wg.Wait()
}
