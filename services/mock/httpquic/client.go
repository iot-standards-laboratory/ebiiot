package httpquic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"services/mock"
	"services/timestamp"
	"services/utils"
	"strings"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
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
	fmt.Println("http quic client start")

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(c.numClients)

	for i := 0; i < c.numClients; i++ {
		go func(id int) {
			defer wg.Done()
			roundTripper := &http3.RoundTripper{
				TLSClientConfig: &tls.Config{
					RootCAs:            pool,
					InsecureSkipVerify: true,
					KeyLogWriter:       nil,
				},
				QuicConfig: &quic.Config{},
				Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, qCfg *quic.Config) (quic.EarlyConnection, error) {
					if strings.Compare("quic.localhost:443", addr) == 0 {
						return quic.DialAddrEarly(c.spAdr, tlsCfg, qCfg)
					}
					return quic.DialAddrEarly(addr, tlsCfg, qCfg)
				},
			}

			defer roundTripper.Close()
			hclient := &http.Client{
				Transport: roundTripper,
			}

			for i := 0; i < c.numMessages; i++ {
				start := time.Now()

				rsp, err := utils.HttpRequest(hclient, int32(c.sizeMessage))
				if err != nil {
					fmt.Println(err)
					return
				}

				body, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(len(body))
				// fmt.Printf("client[%d] - sent %d's message\n", id, i)
				timestamp.Cummulate(int64(time.Since(start).Milliseconds()), timestamp.QUIC)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")
	return nil
}
