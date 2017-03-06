package ads_demo

type ShadowOp struct {
	Aid    int
	Amount float64
	Color  COLOR

	ServerId int
	Depend   VectorClock
}

func (shadow *ShadowOp) apply(bank *BankStorage) {
	account := bank.GetAccount(shadow.Aid)
	account.SetBalance(account.GetBalance() + shadow.Amount)
}
