package timestamp

import (
	"fmt"
	"os"
	"sync"
)

type stampType int

const (
	TCP  stampType = 0
	QUIC stampType = 1
)

var tcpLatency = int64(0)
var tcpCnt = 0
var quicLatency = int64(0)
var quicCnt = 0
var mutex sync.Mutex

func Cummulate(sec int64, st stampType) {
	mutex.Lock()
	defer mutex.Unlock()

	if st == TCP {
		tcpLatency += sec
		tcpCnt++
	} else {
		quicLatency += sec
		quicCnt++
	}

}

func Result() {
	f, err := os.Create("atd.out")
	if err != nil {
		panic(err)
	}

	if tcpCnt != 0 {
		fmt.Fprint(f, tcpLatency/int64(tcpCnt))
	} else {
		fmt.Fprint(f, "0")
	}
	fmt.Fprint(f, " ")

	if quicCnt != 0 {
		fmt.Fprintln(f, quicLatency/int64(quicCnt))
	} else {
		fmt.Fprintln(f, "0")
	}
}
