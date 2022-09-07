package httphybrid

import (
	"fmt"
	"services/mock"
	"services/utils"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	fmt.Println("http hybrid server run")
	go func() {
		err := listenTCP()
		if err != nil {
			panic(err)
		}
	}()

	err := listenQUIC()
	if err != nil {
		return err
	}

	return nil
}

func listenTCP() error {
	fmt.Println("http tcp server start")
	return utils.NewGinMux().RunTLS(":8080", "./assets/cert.pem", "./assets/priv.key")
}

func listenQUIC() error {
	fmt.Println("http quic server start")
	quicConf := &quic.Config{}

	server := http3.Server{
		Handler:    utils.NewMux(),
		Addr:       ":8080",
		QuicConfig: quicConf,
	}

	return server.ListenAndServeTLS("./assets/cert.pem", "./assets/priv.key")
}
