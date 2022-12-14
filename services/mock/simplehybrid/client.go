package simplehybrid

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"services/mock"
	"services/utils"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
)

const tcpRate = 0.4

type Clients struct {
	spAdr          string
	numTCPClients  int
	numQUICClients int
	numTrials      int
	numObjs        int
	sizeMessage    int
}

func NewClients(spAdr string, numClients, numTrials, numObjs, sizeMessage int) mock.Entity {
	numTCPClients := int(float64(numClients) * tcpRate)
	numQUICClients := numClients - numTCPClients
	log.Println("numTCPClients:", numTCPClients)
	log.Println("numQUICClients:", numQUICClients)
	return &Clients{
		spAdr,
		numTCPClients,
		numQUICClients,
		numTrials,
		numObjs,
		sizeMessage,
	}
}

func exchangeTCP(adr string, id, trial, objs, size int) {

}

func exchangeQUIC(adr string, id, trial, objs, size int) {

}

func (c *Clients) Run() error {
	fmt.Println("start hybrid client")

	var wg sync.WaitGroup
	wg.Add(c.numQUICClients + c.numTCPClients)

	for i := 0; i < c.numTCPClients; i++ {
		go func(id, size int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", c.spAdr)
			if err != nil {
				return
			}
			defer conn.Close()

			for i := 0; i < c.numTrials; i++ {
				msg := mock.NewMessage(id, size)
				mock.WritePayload(conn, msg)
				fmt.Printf("tcp client[%d] - sent %d's message\n", id, i)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i, c.sizeMessage)
	}

	tlsConf, err := getTlsConf()
	if err != nil {
		return err
	}
	conn, err := quic.DialAddr(c.spAdr, tlsConf, nil)
	if err != nil {
		return err
	}
	for i := 0; i < c.numQUICClients; i++ {
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
				fmt.Printf("quic client[%d] - sent %d's message\n", id, i)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i, c.sizeMessage)
	}

	wg.Wait()

	fmt.Println("the operation of clients is done!!")

	return nil
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
