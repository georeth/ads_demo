package ads_demo

import "fmt"

type REQ int

const INTEREST_RATE = 0.04
const SERVER_DELAY = 200 // ms

const (
	DEPOSIT REQ = iota
	WITHDRAW
	INTEREST
)

type Request struct {
	Aid    int
	Op     REQ
	Amount float64
}

type Response struct {
	Status  int
	Balance float64
	Message string
}

func (res Response) Print() {
	fmt.Printf("status %d, balance %.2f", res.Status, res.Balance)
	if res.Message != "" {
		fmt.Printf(", message: %s", res.Message)
	}
	fmt.Printf("\n")
}
