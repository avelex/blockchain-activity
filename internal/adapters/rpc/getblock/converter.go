package getblock

import (
	"encoding/json"

	"github.com/avelex/blockchain-activity/internal/types"
	"github.com/avelex/blockchain-activity/pkg/jsonrpc"
)

func blockFromRPC(response jsonrpc.Response) (types.Block, error) {
	encoded, err := json.Marshal(&response.Result)
	if err != nil {
		return types.Block{}, err
	}

	var block types.Block
	if err := json.Unmarshal(encoded, &block); err != nil {
		return types.Block{}, err
	}

	return block, nil
}

func txFromRPC(response jsonrpc.Response) (types.Transaction, error) {
	encoded, err := json.Marshal(&response.Result)
	if err != nil {
		return types.Transaction{}, err
	}

	var tx types.Transaction
	if err := json.Unmarshal(encoded, &tx); err != nil {
		return types.Transaction{}, err
	}

	return tx, nil
}
