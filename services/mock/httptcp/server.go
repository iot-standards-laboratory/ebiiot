package httptcp

import (
	"fmt"
	"net/http"
	"services/mock"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	fmt.Println("http tcp server start")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("Hello world"))
	})
	return http.ListenAndServe(":8080", mux)
}
