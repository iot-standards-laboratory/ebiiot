package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	mock "services/mock"
	"services/mock/echotcp"
	"services/timestamp"
	"syscall"
)

func main() {
	isServer := flag.Bool("server", false, "Is server?")
	exp := flag.String("exp", "quic", "type of experimentation")
	numClients := flag.Int("clients", 10, "number of clients")
	numMessages := flag.Int("messages", 10, "number of messages")
	sizeMessage := flag.Int("size", 100, "size of messages")
	_ = *numClients
	flag.Parse()

	if *isServer {
		runServer(*exp)
		return
	}

	spAdr := flag.Arg(0)
	if len(spAdr) == 0 {
		log.Fatalln("invalid server address")
	}

	var clients mock.Entity
	clients = echotcp.NewClients(spAdr, *numClients, *numMessages, *sizeMessage)

	clients.Run()

	log.Println("done!!")
}

func runServer(exp string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-interrupt
		timestamp.Result()
		os.Exit(0)
	}()

	var s mock.Entity
	s = echotcp.NewServer()

	log.Fatalln(s.Run())
	// run server
}
