package erc20

import (
	"errors"
	"math/big"
	"strings"
)

const transferMethodID = "a9059cbb" // keccak(transfer(address,uint256))

const (
	methodIDLength   = 8
	inputTopicLength = 64
	addressLength    = 40
)

func IsTransfer(txInput string) bool {
	if txInput == "" {
		return false
	}

	input := strings.TrimPrefix(txInput, "0x")
	if len(input) < methodIDLength {
		return false
	}

	return input[:methodIDLength] == transferMethodID
}

func TransferFromTransactionInput(txInput, sender string) (Transfer, error) {
	if !IsTransfer(txInput) {
		return Transfer{}, errors.New("tx not erc20 transfer")
	}

	input := strings.TrimPrefix(txInput, "0x")
	receiverTopic := input[methodIDLength : methodIDLength+inputTopicLength]
	amountTopic := input[len(input)-inputTopicLength:]

	receiver := "0x" + receiverTopic[len(receiverTopic)-addressLength:]
	amountInt, _ := new(big.Int).SetString(amountTopic, 16)

	return Transfer{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amountInt,
	}, nil
}
