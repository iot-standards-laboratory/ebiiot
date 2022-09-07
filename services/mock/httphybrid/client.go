package httphybrid

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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

const tcpRate = 0.4

type Clients struct {
	spAdr          string
	numTCPClients  int
	numQUICClients int
	numMessages    int
	sizeMessage    int
}

func NewClients(spAdr string, numClients, numMessages, sizeMessage int) mock.Entity {
	numTCPClients := int(float64(numClients) * tcpRate)
	numQUICClients := numClients - numTCPClients
	log.Println("numTCPClients:", numTCPClients)
	log.Println("numQUICClients:", numQUICClients)
	return &Clients{
		spAdr,
		numTCPClients,
		numQUICClients,
		numMessages,
		sizeMessage,
	}
}

func (c *Clients) Run() error {
	fmt.Println("http hybrid client start")
	var wg sync.WaitGroup
	wg.Add(c.numQUICClients + c.numTCPClients)

	go c.runQUICClients(&wg)
	go c.runTCPClients(&wg)

	wg.Wait()

	fmt.Println("the operation of clients is done!!")
	return nil
}

func (c *Clients) runQUICClients(wg *sync.WaitGroup) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < c.numQUICClients; i++ {
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

				timestamp.Cummulate(int64(time.Since(start).Milliseconds()), timestamp.QUIC)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i)
	}
}

func (c *Clients) runTCPClients(wg *sync.WaitGroup) {
	for i := 0; i < c.numTCPClients; i++ {
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
				rsp, err := utils.HttpRequest(http.DefaultClient, int32(c.sizeMessage))
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				// fmt.Printf("client[%d] - sent %d's message\n", id, i)
				body, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(len(body))
				timestamp.Cummulate(int64(time.Since(start).Milliseconds()), timestamp.TCP)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i)
	}
}
