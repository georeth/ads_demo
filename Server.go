package ads_demo

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type COLOR int

const (
	BLUE COLOR = iota
	RED
)

type ServerConfig struct {
	Index   int
	Address []string
}

type requestItem struct {
	req     *Request
	resChan chan<- *Response
}

type Server struct {
	id       int
	bank     *BankStorage
	now      *VectorClock
	maxR     int
	hasToken bool

	peers [](*Client) // peer connections

	// events
	tokenChan  chan int         // send or receive token
	reqChan    chan requestItem // client request to process
	shadowChan chan *ShadowOp   // server shadow op to process
	opList     list.List        // received shadows waiting for dependency (ShadowOp)
	redList    list.List        // operations waiting for token to issue (requestItem)

}

func RunServer(config ServerConfig) {
	numServer := len(config.Address)

	server := new(Server)
	server.id = config.Index
	server.bank = NewBankStorage()
	server.now = NewVectorClock(numServer)
	server.tokenChan = make(chan int)
	server.reqChan = make(chan requestItem)
	server.shadowChan = make(chan *ShadowOp)
	server.redList.Init()
	server.opList.Init()

	// setup rpc server
	rpc.Register(server)
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", config.Address[config.Index])
	if e != nil {
		log.Fatalf("server %d: cannot listen at %s\n", config.Index,
			config.Address[config.Index])
	}
	log.Printf("server %d: listen at %s\n", config.Index,
		config.Address[config.Index])

	// setup peer connection
	go func() {
		server.peers = make([]*Client, numServer)
		for i := 0; i < numServer; i++ {
			if i == config.Index {
				continue
			}
			server.peers[i] = NewClient(config.Address[i])
		}
		log.Printf("server %d: peer connection established\n", config.Index)
		server.mainLoop()
	}()

	http.Serve(listener, nil)
}

func (server *Server) setTokenTimeout() {
	go func() {
		time.Sleep(time.Second)
		server.tokenChan <- 0
	}()
}

func (server *Server) primary() bool {
	return server.hasToken && server.maxR == server.now.Red()
}

func (server *Server) generateShadow(req *Request, primary bool) (shadow *ShadowOp, res *Response, ok bool) {
	shadow = &ShadowOp{Aid: req.Aid, Depend: server.now.Copy(), ServerId: server.id}
	// fill amount, color
	balance := server.bank.GetAccount(req.Aid).GetBalance()

	switch req.Op {
	case DEPOSIT:
		shadow.Amount = req.Amount
		shadow.Color = BLUE
		res = &Response{Status: 0, Balance: balance + req.Amount}
		ok = true
	case WITHDRAW:
		// Gemini use OCC
		// (generate shadow op, then verify when hasToken. client retry if failed)
		// but for this application, delay until issue
		if primary {
			if balance >= req.Amount {
				shadow.Amount = -req.Amount
				shadow.Color = RED
				res = &Response{Status: 0, Balance: balance - req.Amount}
			} else {
				shadow.Amount = 0
				shadow.Color = BLUE
				res = &Response{Status: -1, Balance: balance,
					Message: "Insufficient balance"}
			}
			ok = true
		} else {
			ok = false
		}
	case INTEREST:
		if primary {
			delta := server.bank.GetAccount(req.Aid).ComputeInterest()
			shadow.Amount = delta
			shadow.Color = BLUE
			res = &Response{Status: 0, Balance: balance + delta}
			ok = true

			shadow.Color = RED
		} else {
			ok = false
		}
	default:
		panic("Unknown operation")
	}
	return shadow, res, ok
}

func (server *Server) doRequest(reqItem *requestItem) bool {
	primary := server.primary()
	req := reqItem.req

	// verify request
	if req.Aid < 0 || req.Aid >= NumAccounts {
		reqItem.resChan <- &Response{Status: -1, Message: "Invalid Account Id"}
		return true
	}

	// try generate shadow op
	shadow, res, ok := server.generateShadow(req, primary)
	if ok {
		reqItem.resChan <- res
		server.dispatchShadowOp(shadow)
		return true
	}
	if !ok && primary {
		fmt.Printf("failed %d: \n", req.Op)
	}
	return false
}

func (server *Server) dispatchShadowOp(shadow *ShadowOp) {
	shadow.apply(server.bank)
	server.now.Tick(shadow.ServerId, shadow.Color)
	if server.now.Red() > server.maxR {
		server.maxR = server.now.Red()
	}

	for _, p := range server.peers {
		if p != nil {
			p.AddShadowOpAsync(shadow)
		}
	}
}

func (server *Server) mainLoop() {
	if server.id == 0 {
		server.setTokenTimeout()
		server.hasToken = true
	}
	// now process all events in a single-thread loop
	for {
		select {
		// send or receive token
		case maxR := <-server.tokenChan:
			server.now.Print(server.id)
			if server.hasToken {
				server.hasToken = false
				nextId := (server.id + 1) % len(server.peers)
				server.peers[nextId].PassToken(server.maxR)
			} else {
				server.maxR = maxR
				server.hasToken = true
				server.setTokenTimeout()
			}
		case shadow := <-server.shadowChan:
			server.opList.PushBack(shadow)
		case reqItem := <-server.reqChan:
			ok := server.doRequest(&reqItem)
			if !ok {
				server.redList.PushBack(&reqItem)
			}
		}

		// process opList
		for todo := true; todo; {
			todo = false
			for e := server.opList.Front(); e != nil; {
				shadow := e.Value.(*ShadowOp)
				if shadow.Depend.Ready(server.now) {
					shadow.apply(server.bank)
					server.now.Tick(shadow.ServerId, shadow.Color)
					if server.now.Red() > server.maxR {
						server.maxR = server.now.Red()
					}

					todo = true

					var cur *list.Element
					cur, e = e, e.Next()
					server.opList.Remove(cur)
				} else {
					e = e.Next()
				}
			}
		}

		// process redList
		if server.primary() {
			for e := server.redList.Front(); e != nil; e = e.Next() {
				reqItem := e.Value.(*requestItem)
				ok := server.doRequest(reqItem)

				if !ok {
					log.Fatalf("server %d: process redList fail\n", server.id)
				}
			}
			server.redList.Init()
		}
	}
}

/***************** RPC Part ******************/
/*****************  Server  ******************/

func (server *Server) PassToken(maxR int, _ *struct{}) error {
	server.tokenChan <- maxR
	return nil
}

func (server *Server) AddShadowOp(shadow ShadowOp, _ *struct{}) error {
	server.shadowChan <- &shadow
	return nil
}

/*****************  Client  ******************/

func (server *Server) Request(req Request, res *Response) error {
	// convert sync rpc handling to async by channel
	var reqItem requestItem
	resChan := make(chan *Response)

	reqItem.req = &req
	reqItem.resChan = resChan
	server.reqChan <- reqItem
	resp, ok := <-resChan

	*res = *resp

	if !ok {
		return errors.New("Server.Request Fail")
	}
	return nil
}

/*****************  Misc  ******************/
func (server *Server) Dump(t1 int, t2 *int) error {
	fmt.Printf("server %d\n", server.id)
	return nil
}
