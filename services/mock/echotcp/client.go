package echotcp

import (
	"log"
	"net"
	"services/mock"
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
	var wg sync.WaitGroup
	wg.Add(c.numClients)

	for i := 0; i < c.numClients; i++ {
		go func(id, size int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", c.spAdr)
			if err != nil {
				return
			}
			defer conn.Close()

			for i := 0; i < c.numMessages; i++ {
				msg := mock.NewMessage(id, size)
				mock.WritePayload(conn, msg)
				log.Printf("client[%d] - sent %d's message\n", id, i)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i, c.sizeMessage)
	}

	wg.Wait()
	log.Println("the operation of clients is done!!")

	return nil
}
