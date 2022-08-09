package timestamp

import (
	"fmt"
	"os"
	"sync"
)

var latency = int64(0)
var cnt = 0
var mutex sync.Mutex

func Cummulate(sec int64) {
	mutex.Lock()
	defer mutex.Unlock()
	latency += sec
	cnt++
}

func Result() {
	f, err := os.Create("atd.out")
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(f, latency/int64(cnt))
}
