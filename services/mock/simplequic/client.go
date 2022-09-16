package simplequic

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"services/mock"
	"services/utils"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
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

func exchange(adr string, id, trials, objs, size int) error {
	keylog, err := os.Create(fmt.Sprintf("key-%d.log", id))
	if err != nil {
		return err
	}
	defer keylog.Close()

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		KeyLogWriter:       keylog,
		NextProtos:         []string{"quic-echo-example"},
	}

	if err != nil {
		return err
	}
	conn, err := quic.DialAddr(adr, tlsConf, nil)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(objs)

	for i := 0; i < objs; i++ {
		go func() {
			defer wg.Done()
			stream, err := conn.OpenStreamSync(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
			defer stream.Close()

			for i := 0; i < trials; i++ {
				msg := mock.NewMessage(id, size)
				mock.WritePayload(stream, msg)
				fmt.Printf("client[%d] - sent %d's message\n", id, i)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	return nil
}

func (c *Clients) Run() error {
	fmt.Println("start quic client")
	var wg sync.WaitGroup
	wg.Add(c.numClients)

	for i := 0; i < c.numClients; i++ {
		go func(adr string, id, trials, objs, size int) {
			defer wg.Done()
			exchange(c.spAdr, id, trials, objs, size)
		}(c.spAdr, i, c.numTrials, c.numObjs, c.sizeMessage)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")

	return nil
}
