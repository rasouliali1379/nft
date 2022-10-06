package talan

type Response struct {
	StatusCode int    `json:"statusCode"`
	Data       any    `json:"data"`
	Error      string `json:"error"`
}

type GeneratedAddressDto struct {
	Mnemonic      string `json:"mnemonic"`
	PublicAddress string `json:"publicAddress"`
	PrivateKey    string `json:"privateKey"`
}

type TransactionDto struct {
	TxId          string  `json:"txId"`
	BlockHeight   int64   `json:"blockHeight"`
	Timestamp     int64   `json:"timestamp"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	Type          string  `json:"type"`
}
