package httptcp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
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
			pool, err := x509.SystemCertPool()
			if err != nil {
				log.Fatal(err)
			}

			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				// DualStack: true,
			}
			payload := make([]byte, c.sizeMessage)
			for i := 0; i < c.sizeMessage; i++ {
				payload[i] = 'b'
			}

			http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				if addr == "quic.localhost:443" {
					addr = c.spAdr
				}
				return dialer.DialContext(ctx, network, addr)
			}

			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: false,
			}

			for i := 0; i < c.numMessages; i++ {
				start := time.Now()
				rsp, err := http.Post("https://quic.localhost", "text/plain", bytes.NewReader(payload))
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				// fmt.Printf("client[%d] - sent %d's message\n", id, i)
				fmt.Println(rsp.Proto)
				timestamp.Cummulate(int64(time.Since(start).Milliseconds()), timestamp.TCP)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")
	return nil
}
