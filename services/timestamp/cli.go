package timestamp

import (
	"fmt"
	"os"
	"services/utils"
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

	utils.Done()
}

func Result() {
	f, err := os.OpenFile("atd.out", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	if tcpCnt != 0 {
		fmt.Fprintf(f, "%d %d\n", tcpCnt, tcpLatency/int64(tcpCnt))
	} else {
		fmt.Fprintln(f, "0 0")
	}

	if quicCnt != 0 {
		fmt.Fprintf(f, "%d %d\n", quicCnt, quicLatency/int64(quicCnt))
	} else {
		fmt.Fprintln(f, "0 0")
	}
}
