package ads_demo

import "net/rpc"
import "log"
import "time"

type Client struct {
	rpcClient *rpc.Client
}

func NewClient(addr string) *Client {
	var c Client
	for c.rpcClient == nil {
		var err error
		c.rpcClient, err = rpc.DialHTTP("tcp", addr)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
		}
	}
	return &c
}

func (c *Client) PassToken(maxR int) {
	var empty struct{}
	go func() {
		time.Sleep(SERVER_DELAY * time.Millisecond)
		err := c.rpcClient.Call("Server.PassToken", maxR, &empty)
		if err != nil {
			log.Fatalf("client.Dump() : %s\n", err.Error())
		}
	}()
}

func (c *Client) AddShadowOpAsync(op *ShadowOp) {
	var empty struct{}
	go func() {
		time.Sleep(SERVER_DELAY * time.Millisecond)
		err := c.rpcClient.Call("Server.AddShadowOp", op, &empty)
		// FIXME
		if err != nil {
			log.Fatalf("client.AddShadowOp() : %s\n", err.Error())
		}
	}()
}

func (c *Client) Request(req Request) Response {
	var res Response
	err := c.rpcClient.Call("Server.Request", req, &res)

	if err != nil {
		log.Fatalf("client.Request() : %s\n", err.Error())
	}
	return res
}

func (c *Client) Dump() {
	var res int
	err := c.rpcClient.Call("Server.Dump", 1, &res)
	if err != nil {
		log.Fatalf("client.Dump() : %s\n", err.Error())
	}
}
