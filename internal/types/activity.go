package types

type BlockchainActivity struct {
	BlockNumber int64             `json:"block_number"`
	Chain       int64             `json:"chain"`
	Addresses   []AddressActivity `json:"addresses"`
}

type AddressActivity struct {
	Address  string `json:"address"`
	Activity uint   `json:"activity"`
}
