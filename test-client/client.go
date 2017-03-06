package main

import (
	"ads_demo"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const CLIENT_PER_SERVER = 100
const RUNTIME = 10

type perf struct {
	op      [3]int64
	latency [3]int64
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: ./test-client addr1 addr2 ...\n")
	panic("Wrong command line argument")
}

func randint(max int) int {
	return rand.Int() % max
}
func randop() ads_demo.REQ {
	r := randint(100)
	// 40 : 40 : 20
	if r < 40 {
		return ads_demo.DEPOSIT
	} else if r < 80 {
		return ads_demo.WITHDRAW
	} else {
		return ads_demo.INTEREST
	}
}

func test(client []*ads_demo.Client, p *perf, endChan <-chan bool, num_server int) {
	for {
		select {
		case <-endChan:
			return
		default:
			start := time.Now()
			req := ads_demo.Request{
				Aid:    randint(ads_demo.NumAccounts),
				Amount: float64(randint(200)),
				Op:     randop()}

			client[randint(num_server)].Request(req)
			latency := int64(time.Now().Sub(start)) / int64(time.Microsecond)

			p.op[req.Op]++
			p.latency[req.Op] += latency
		}
	}
}

func main() {
	args := os.Args[1:]
	num_server := len(args)

	fmt.Printf("num_server : %d\n", num_server)

	if len(args) < 1 {
		Usage()
	}

	perfs := make([]perf, CLIENT_PER_SERVER)
	endChan := make(chan bool)

	for i := 0; i < CLIENT_PER_SERVER; i++ {
		clients := make([]*ads_demo.Client, num_server)
		for j := 0; j < num_server; j++ {
			clients[j] = ads_demo.NewClient(args[j])
			fmt.Printf("connected\n")
		}
		go test(clients, &perfs[i], endChan, num_server)
	}

	time.Sleep(RUNTIME * time.Second)
	for k := 0; k < CLIENT_PER_SERVER; k++ {
		endChan <- false
	}

	// output summary
	var sum perf
	for i := 0; i < CLIENT_PER_SERVER; i++ {
		sum.op[ads_demo.DEPOSIT] += perfs[i].op[ads_demo.DEPOSIT]
		sum.latency[ads_demo.DEPOSIT] += perfs[i].latency[ads_demo.DEPOSIT]

		sum.op[ads_demo.WITHDRAW] += perfs[i].op[ads_demo.WITHDRAW]
		sum.latency[ads_demo.WITHDRAW] += perfs[i].latency[ads_demo.WITHDRAW]

		sum.op[ads_demo.INTEREST] += perfs[i].op[ads_demo.INTEREST]
		sum.latency[ads_demo.INTEREST] += perfs[i].latency[ads_demo.INTEREST]
	}

	fmt.Printf("latency deposit:%dus, withdraw %dus, interest %dus\n",
		sum.latency[ads_demo.DEPOSIT]/sum.op[ads_demo.DEPOSIT],
		sum.latency[ads_demo.WITHDRAW]/sum.op[ads_demo.WITHDRAW],
		sum.latency[ads_demo.INTEREST]/sum.op[ads_demo.INTEREST],
	)
}
