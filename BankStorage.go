package ads_demo

const NumAccounts = 10000
const InitBalance = 1000.0

type BankStorage struct {
	accounts map[int]*Account
}

func (bank *BankStorage) GetAccount(id int) *Account {
	return bank.accounts[id]
}

func NewBankStorage() *BankStorage {
	var bank BankStorage
	bank.accounts = make(map[int]*Account)
	for i := 0; i < NumAccounts; i++ {
		bank.accounts[i] = &Account{i, InitBalance}
	}
	return &bank
}
