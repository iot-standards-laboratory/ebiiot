package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"os"
	"os/signal"
	mock "services/mock"
	"services/mock/httphybrid"
	"services/mock/httpquic"
	"services/mock/httptcp"
	"services/mock/simplehybrid"
	"services/mock/simplequic"
	"services/mock/simpletcp"
	"services/timestamp"
	"services/utils"
	"strconv"
	"strings"
	"syscall"
)

func getHashValue(b string) uint32 {
	return crc32.Checksum([]byte(b), crc32.MakeTable(crc32.IEEE))
}

var clientGenerators = map[uint32]map[uint32]func(string, int, int, int, int) mock.Entity{}
var serverGenerators = map[uint32]map[uint32]func() mock.Entity{}

func init() {
	clientSimpleGenerators := map[uint32]func(string, int, int, int, int) mock.Entity{}
	serverSimpleGenerators := map[uint32]func() mock.Entity{}
	clientSimpleGenerators[getHashValue("tcp")] = simpletcp.NewClients
	serverSimpleGenerators[getHashValue("tcp")] = simpletcp.NewServer
	clientSimpleGenerators[getHashValue("quic")] = simplequic.NewClients
	serverSimpleGenerators[getHashValue("quic")] = simplequic.NewServer
	clientSimpleGenerators[getHashValue("hybrid")] = simplehybrid.NewClients
	serverSimpleGenerators[getHashValue("hybrid")] = simplehybrid.NewServer

	clientHttpGenerators := map[uint32]func(string, int, int, int, int) mock.Entity{}
	serverHttpGenerators := map[uint32]func() mock.Entity{}

	clientHttpGenerators[getHashValue("tcp")] = httptcp.NewClients
	serverHttpGenerators[getHashValue("tcp")] = httptcp.NewServer
	clientHttpGenerators[getHashValue("quic")] = httpquic.NewClients
	serverHttpGenerators[getHashValue("quic")] = httpquic.NewServer
	clientHttpGenerators[getHashValue("hybrid")] = httphybrid.NewClients
	serverHttpGenerators[getHashValue("hybrid")] = httphybrid.NewServer

	clientGenerators[getHashValue("simple")] = clientSimpleGenerators
	serverGenerators[getHashValue("simple")] = serverSimpleGenerators
	clientGenerators[getHashValue("http")] = clientHttpGenerators
	serverGenerators[getHashValue("http")] = serverHttpGenerators
}

func main() {
	exp := flag.String("exp", "simple", "type of experimentation")
	isServer := flag.Bool("server", false, "Is server?")
	proto := flag.String("proto", "quic", "transport protocol")
	numClients := flag.Int("clients", 10, "number of clients")
	numTrials := flag.Int("trials", 10, "number of trials per client")
	numObjs := flag.Int("objs", 10, "number of objs")
	sizeMessage := flag.Int("size", 100, "size of messages")
	_ = *numClients
	_ = *exp
	flag.Parse()

	if *isServer {
		go runCounter()
		runServer(*exp, *proto)

		return
	}

	spAdr := flag.Arg(0)
	if len(spAdr) == 0 {
		log.Fatalln("invalid server address")
	}

	var clients = clientGenerators[getHashValue(*exp)][getHashValue(*proto)](spAdr, *numClients, *numTrials, *numObjs, *sizeMessage)
	go clients.Run()
	http.Get(fmt.Sprintf("http://localhost:10000/%d", (*numClients)*(*numObjs)*(*numTrials)))
	fmt.Println("done!!")
	if strings.Compare(*exp, "http") == 0 {
		timestamp.Result()
	}
}

func runCounter() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		count, _ := strconv.ParseInt(req.URL.Path[1:], 0, 32)
		fmt.Println(count)
		utils.Wait(int(count))
		resp.WriteHeader(http.StatusOK)
	})

	log.Fatalln(http.ListenAndServe(":10000", mux))
}

func runServer(exp, proto string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)

	go func() {
		log.Println(<-interrupt)
		if strings.Compare(exp, "simple") == 0 {
			timestamp.Result()
		}
		os.Exit(0)
	}()

	s := serverGenerators[getHashValue(exp)][getHashValue(proto)]()

	// run server
	fmt.Println("server run:", s.Run())
}
