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
func getTlsConf() (*tls.Config, error) {
	keylog, err := os.Create("./ssl-key.log")
	if err != nil {
		return nil, err
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		KeyLogWriter:       keylog,
		NextProtos:         []string{"quic-echo-example"},
	}

	return tlsConf, nil
}

func (c *Clients) Run() error {
	fmt.Println("start quic client")
	var wg sync.WaitGroup
	wg.Add(c.numClients)

	tlsConf, err := getTlsConf()
	if err != nil {
		return err
	}
	conn, err := quic.DialAddr(c.spAdr, tlsConf, nil)
	if err != nil {
		return err
	}
	for i := 0; i < c.numClients; i++ {
		go func(id, size int) {
			defer wg.Done()
			stream, err := conn.OpenStreamSync(context.Background())
			if err != nil {
				return
			}
			defer stream.Close()

			for i := 0; i < c.numTrials; i++ {
				msg := mock.NewMessage(id, size)
				mock.WritePayload(stream, msg)
				fmt.Printf("client[%d] - sent %d's message\n", id, i)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i, c.sizeMessage)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")

	return nil
}
