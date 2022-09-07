package httptcp

import (
	"fmt"
	"net/http"
	"services/mock"
	"services/utils"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	fmt.Println("http tcp server start")
	return http.ListenAndServeTLS(":8080", "./assets/cert.pem", "./assets/priv.key", utils.NewMux())
	// return utils.NewGinMux().RunTLS(":8080", "./assets/cert.pem", "./assets/priv.key")
}
