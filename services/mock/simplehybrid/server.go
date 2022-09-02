package simplehybrid

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"services/mock"
	"services/timestamp"
	"sync"

	"github.com/lucas-clemente/quic-go"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	fmt.Println("hybrid run server")

	// ctx, _ := context.WithCancel(context.Background())
	ctx := context.Background()
	err := listenTCP(ctx)
	if err != nil {
		return err
	}

	return nil
}

func listenTCP(ctx context.Context) error {
	listener, err := net.Listen(
		"tcp",
		":8080",
	)
	if err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()
	terminated := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			payload, err := mock.ReadPayload(conn)
			if err != nil {
				terminated <- struct{}{}
				return
			}
			msg := mock.ParseMsg(payload)
			timestamp.Cummulate(msg.Latency(), timestamp.TCP)
		}
	}()

	go func() {
		wg.Wait()
		terminated <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-terminated:
		log.Println("tt: terminated")
		return nil
	}
}

func listenQUIC(ctx context.Context) error {
	listener, err := quic.ListenAddrEarly(
		":8080",
		generateTLSConfig(),
		nil,
	)
	if err != nil {
		return err
	}

	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}

	terminated := make(chan interface{})

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return err
		}
		go func() {
			for {
				payload, err := mock.ReadPayload(stream)
				if err != nil {
					terminated <- struct{}{}
					return
				}

				msg := mock.ParseMsg(payload)
				timestamp.Cummulate(msg.Latency(), timestamp.QUIC)
			}
		}()
	}
	// select {
	// case <-ctx.Done():
	// 	return nil
	// case <-terminated:
	// 	return nil
	// }

}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
