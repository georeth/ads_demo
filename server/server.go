package main

import "os"
import "fmt"
import "strconv"
import "ads_demo"

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: ./server index addr1 addr2 ...\n")
	panic("Wrong command line argument")
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		Usage()
	}
	index, e := strconv.ParseInt(args[0], 0, 16)
	if e != nil {
		Usage()
	}
	addr := args[1:]

	ads_demo.RunServer(ads_demo.ServerConfig{int(index), addr})
}
