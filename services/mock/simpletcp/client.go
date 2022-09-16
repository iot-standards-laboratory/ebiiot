package simpletcp

import (
	"fmt"
	"net"
	"services/mock"
	"services/utils"
	"sync"
	"time"
)

type Clients struct {
	spAdr       string
	numClients  int
	numTrials   int
	numObjs     int
	sizeMessage int
}

func NewClients(spAdr string, numClients, numTrials, numObjs, sizeMessage int) mock.Entity {
	return &Clients{
		spAdr,
		numClients,
		numTrials,
		numObjs,
		sizeMessage,
	}
}

func (c *Clients) Run() error {
	fmt.Println("start tcp client")
	var wg sync.WaitGroup
	wg.Add(c.numClients)
	for i := 0; i < c.numClients; i++ {
		go func(id, size int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", c.spAdr)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			for i := 0; i < c.numTrials; i++ {
				msg := mock.NewMessage(id, size)
				mock.WritePayload(conn, msg)
				fmt.Printf("client[%d] - sent %d's message\n", id, i)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i, c.sizeMessage)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")

	return nil
}
