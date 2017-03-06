package ads_demo

type Account struct {
	id      int
	balance float64
}

func (acc *Account) GetId() int {
	return acc.id
}
func (acc *Account) GetBalance() float64 {
	return acc.balance
}
func (acc *Account) ComputeInterest() float64 {
	return acc.balance * INTEREST_RATE
}
func (acc *Account) SetBalance(balance float64) {
	acc.balance = balance
}
