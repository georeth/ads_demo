package main

import (
	"ads_demo"
	_ "fmt"
	_ "os"
	"time"
)

const CLIENT_PER_SERVER = 100
const RUNTIME = 10

func main() {
	clients := []*ads_demo.Client{
		ads_demo.NewClient("localhost:13000"),
		ads_demo.NewClient("localhost:13001"),
		ads_demo.NewClient("localhost:13002")}

	clients[0].Request(ads_demo.Request{
		Aid:    20,
		Op:     ads_demo.DEPOSIT,
		Amount: 1000}).Print()
	clients[1].Request(ads_demo.Request{
		Aid:    20,
		Op:     ads_demo.DEPOSIT,
		Amount: 1100}).Print()
	clients[0].Request(ads_demo.Request{
		Aid: 20,
		Op:  ads_demo.INTEREST}).Print()

	time.Sleep(time.Second)

	clients[1].Request(ads_demo.Request{
		Aid:    20,
		Op:     ads_demo.WITHDRAW,
		Amount: 2500}).Print()

}
