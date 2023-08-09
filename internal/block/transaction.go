package block

type TXInput struct {
	Txid      []byte // ID of tx this input refers
	VoutIndex int    // index of an output in the tx
	Signature string
}

type TXOutput struct {
	// stores the number of satoshis, which is 0.00000001 BTC.
	// this is the smallest unit of currency in Bitcoin
	Value int64
	// we are not implementing the whole srcipt thing here, so just pubkey
	PubKeyHash string
}

type Transaction struct {
	ID   []byte // hash of this tx
	Vin  []TXInput
	Vout []TXOutput
}

const subsidy = 10000

// create a new tx, which has an output to reward the miner
// this output
func NewCoinbaseTransaction() *Transaction {
	return nil
}
