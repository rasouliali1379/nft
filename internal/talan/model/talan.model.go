package talan

type Talan struct {
	Balance      float64
	Address      Address
	Transactions []Transaction
}

type Transaction struct {
	ID            string
	BlockHeight   int64
	TimeStamp     int64
	Amount        float64
	Confirmations int64
	Type          TransactionType
}

type TransactionType string

const (
	TransactionTypeReceive TransactionType = "receive"
	TransactionTypeSend    TransactionType = "send"
)

type Address struct {
	Mnemonic      string
	PublicAddress string
	PrivateKey    string
}
