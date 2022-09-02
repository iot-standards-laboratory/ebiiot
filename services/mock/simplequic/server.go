package simplequic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"services/mock"
	"services/timestamp"

	"github.com/lucas-clemente/quic-go"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	listener, err := quic.ListenAddrEarly(
		":8080",
		generateTLSConfig(),
		nil,
	)

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}

		go func() {
			for {
				stream, err := conn.AcceptStream(context.Background())
				if err != nil {
					return
				}

				go func() {
					for {
						payload, err := mock.ReadPayload(stream)
						if err != nil {
							return
						}

						msg := mock.ParseMsg(payload)
						timestamp.Cummulate(msg.Latency(), timestamp.QUIC)
					}
				}()
			}
		}()
	}
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
