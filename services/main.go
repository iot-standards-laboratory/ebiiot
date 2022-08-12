package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	mock "services/mock"
	"services/mock/simplequic"
	"services/mock/simpletcp"
	"services/timestamp"
	"strings"
	"syscall"
)

func main() {
	exp := flag.String("exp", "simple", "type of experimentation")
	isServer := flag.Bool("server", false, "Is server?")
	proto := flag.String("proto", "quic", "transport protocol")
	numClients := flag.Int("clients", 10, "number of clients")
	numMessages := flag.Int("messages", 10, "number of messages")
	sizeMessage := flag.Int("size", 100, "size of messages")
	_ = *numClients
	_ = *exp
	flag.Parse()

	if *isServer {
		runServer(*proto)
		return
	}

	spAdr := flag.Arg(0)
	if len(spAdr) == 0 {
		log.Fatalln("invalid server address")
	}

	var clients mock.Entity
	if strings.Compare(*proto, "tcp") == 0 {
		clients = simpletcp.NewClients(spAdr, *numClients, *numMessages, *sizeMessage)
	} else if strings.Compare(*proto, "quic") == 0 {
		clients = simplequic.NewClients(spAdr, *numClients, *numMessages, *sizeMessage)
	}

	clients.Run()

	fmt.Println("done!!")
}

func runServer(proto string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-interrupt
		timestamp.Result()
		os.Exit(0)
	}()

	var s mock.Entity
	if strings.Compare(proto, "tcp") == 0 {
		s = simpletcp.NewServer()
	} else if strings.Compare(proto, "quic") == 0 {
		s = simplequic.NewServer()
	}

	log.Fatalln(s.Run())
	// run server
}
