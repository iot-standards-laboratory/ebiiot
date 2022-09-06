package httpquic

import (
	"fmt"
	"net/http"
	"services/mock"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	fmt.Println("http quic server start")
	quicConf := &quic.Config{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("Hello world"))
	})

	server := http3.Server{
		Handler:    mux,
		Addr:       ":8080",
		QuicConfig: quicConf,
	}

	return server.ListenAndServeTLS("./assets/cert.pem", "./assets/priv.key")
}
