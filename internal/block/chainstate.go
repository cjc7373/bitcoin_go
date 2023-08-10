package block

const utxoBucket = "chainstate"

type UTXOSet struct {
	Blockchain *Blockchain
}

// rebuild UXTO set
func (u UTXOSet) Reindex() {

}
