package httptcp

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"services/mock"
	"services/timestamp"
	"services/utils"
	"sync"
	"time"
)

type Clients struct {
	spAdr       string
	numClients  int
	numMessages int
	sizeMessage int
}

func NewClients(spAdr string, numClients, numMessages, sizeMessage int) mock.Entity {
	return &Clients{
		spAdr,
		numClients,
		numMessages,
		sizeMessage,
	}
}

func (c *Clients) Run() error {
	fmt.Println("http tcp client start")

	var wg sync.WaitGroup
	wg.Add(c.numClients)
	for i := 0; i < c.numClients; i++ {
		go func(id int) {
			defer wg.Done()
			payload := make([]byte, c.sizeMessage)
			for i := 0; i < c.sizeMessage; i++ {
				payload[i] = 'b'
			}

			for i := 0; i < c.numMessages; i++ {
				start := time.Now()
				_, err := http.Post(fmt.Sprintf("http://%s", c.spAdr), "text/plain", bytes.NewReader(payload))
				if err != nil {
					log.Println(err)
					return
				}

				fmt.Printf("client[%d] - sent %d's message\n", id, i)
				timestamp.Cummulate(int64(time.Since(start).Milliseconds()), timestamp.TCP)

				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")
	return nil
}
