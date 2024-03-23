package erc20

import (
	"fmt"
	"testing"
)

func TestIsTransfer(t *testing.T) {
	txInput := "0xa9059cbb00000000000000000000000075e89d5979e4f6fba9f97c104c2f0afb3f1dcb8800000000000000000000000000000000000000000000000000000000bea5237b"
	want := true
	got := IsTransfer(txInput)
	if want != got {
		t.Fail()
	}

	tr, err := TransferFromTransactionInput(txInput, "0x5f65f7b609678448494De4C87521CdF6cEf1e932")
	if err != nil {
		t.Fail()
	}

	fmt.Printf("tr: %v\n", tr)
}
