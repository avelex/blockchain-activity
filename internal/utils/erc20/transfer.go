package erc20

import "math/big"

type Transfer struct {
	Sender   string
	Receiver string
	Amount   *big.Int
}
