package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"os/signal"
	mock "services/mock"
	"services/mock/simplequic"
	"services/mock/simpletcp"
	"services/timestamp"
	"syscall"
)

func getHashValue(b string) uint32 {
	return crc32.Checksum([]byte(b), crc32.MakeTable(crc32.IEEE))
}

var clientGenerators = map[uint32]map[uint32]func(string, int, int, int) mock.Entity{}
var serverGenerators = map[uint32]map[uint32]func() mock.Entity{}

func init() {
	clientSimpleGenerators := map[uint32]func(string, int, int, int) mock.Entity{}
	serverSimpleGenerators := map[uint32]func() mock.Entity{}
	clientSimpleGenerators[getHashValue("tcp")] = simpletcp.NewClients
	serverSimpleGenerators[getHashValue("tcp")] = simpletcp.NewServer
	clientSimpleGenerators[getHashValue("quic")] = simplequic.NewClients
	serverSimpleGenerators[getHashValue("quic")] = simplequic.NewServer
	clientGenerators[getHashValue("simple")] = clientSimpleGenerators
	serverGenerators[getHashValue("simple")] = serverSimpleGenerators
}

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
		runServer(*exp, *proto)
		return
	}

	spAdr := flag.Arg(0)
	if len(spAdr) == 0 {
		log.Fatalln("invalid server address")
	}

	var clients = clientGenerators[getHashValue(*exp)][getHashValue(*proto)](spAdr, *numClients, *numMessages, *sizeMessage)
	clients.Run()

	fmt.Println("done!!")
}

func runServer(exp, proto string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-interrupt
		timestamp.Result()
		os.Exit(0)
	}()

	s := serverGenerators[getHashValue(exp)][getHashValue(proto)]()

	log.Fatalln(s.Run())
	// run server
}
