package httpquic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

var pool *x509.CertPool

func init() {
	var err error
	pool, err = x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
}

func exchange(adr string, id, trial, objs, size int) {
	f, err := os.Create(fmt.Sprintf("key-%d-%d.log", id, trial))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       f,
		},
		QuicConfig: &quic.Config{},
		Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, qCfg *quic.Config) (quic.EarlyConnection, error) {
			if strings.Compare("quic.localhost:443", addr) == 0 {
				return quic.DialAddrEarly(adr, tlsCfg, qCfg)
			}
			return quic.DialAddrEarly(addr, tlsCfg, qCfg)
		},
	}

	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	var wg sync.WaitGroup
	wg.Add(objs)
	for i := 0; i < objs; i++ {
		go func(id, trial, obj int) {
			defer wg.Done()
			start := time.Now()

			rsp, err := utils.HttpRequest(hclient, size)
			if err != nil {
				fmt.Println(err)
				return
			}

			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Printf("id: %d, trial: %d, obj: %d - len: %d\n", id, trial, obj, len(body))
			timestamp.Cummulate(int64(time.Since(start).Milliseconds()), timestamp.QUIC)
		}(id, trial, i)
	}

	wg.Wait()
}
func (c *Clients) Run() error {
	fmt.Println("http quic client start")

	var wg sync.WaitGroup
	wg.Add(c.numClients)

	for i := 0; i < c.numClients; i++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < c.numTrials; i++ {
				exchange(c.spAdr, id, i, c.numObjs, c.sizeMessage)
				time.Sleep(time.Duration(utils.GetSleepTime()) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("the operation of clients is done!!")
	return nil
}
