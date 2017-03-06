package main

import "ads_demo"
import "fmt"
import "os"
import "bufio"
import "time"

func main() {
	reader := bufio.NewReader(os.Stdin)
	client := ads_demo.NewClient(os.Args[1])

	fmt.Printf("Connected\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		var cmd string
		var arg1 int     // aid
		var arg2 float64 // amount
		var req ads_demo.Request
		var res ads_demo.Response

		start := time.Now()

		n, _ := fmt.Sscanf(line, "%s%d%f", &cmd, &arg1, &arg2)
		if cmd == "deposit" && n == 3 {
			req = ads_demo.Request{Op: ads_demo.DEPOSIT, Aid: arg1, Amount: arg2}
			res = client.Request(req)
		} else if cmd == "withdraw" && n == 3 {
			req = ads_demo.Request{Op: ads_demo.WITHDRAW, Aid: arg1, Amount: arg2}
			res = client.Request(req)
		} else if cmd == "interest" && n == 2 {
			req = ads_demo.Request{Op: ads_demo.INTEREST, Aid: arg1}
			res = client.Request(req)
		} else {
			fmt.Printf("Retry\n")
			continue
		}

		fmt.Printf("\tstatus: %d, balance: %.2f\n", res.Status, res.Balance)
		if res.Message != "" {
			fmt.Printf("\tmessage: %s\n", res.Message)
		}
		fmt.Printf("\tlatency %.3f ms\n",
			float64(time.Now().Sub(start))/float64(time.Millisecond))
		fmt.Printf("\n")
	}
}
