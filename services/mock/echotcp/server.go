package echotcp

import (
	"net"
	"services/mock"
	"services/timestamp"
)

type Server struct{}

func NewServer() mock.Entity {
	return &Server{}
}

func (s *Server) Run() error {
	listener, err := net.Listen(
		"tcp",
		":8080",
	)

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			defer conn.Close()
			for {
				payload, err := mock.ReadPayload(conn)
				if err != nil {
					return
				}

				msg := mock.ParseMsg(payload)
				timestamp.Cummulate(msg.Latency())
			}
		}()
	}

}
