package getblock

import (
	"context"
	"fmt"
	"strconv"

	"github.com/avelex/blockchain-activity/internal/types"
	"github.com/avelex/blockchain-activity/pkg/jsonrpc"
)

const (
	getBlockBase = "https://go.getblock.io"
	getBlockID   = "getblock.io"
)

const (
	blockNumberMethod          = "eth_blockNumber"
	chainIDMethod              = "eth_chainId"
	getBlockByNumberMethod     = "eth_getBlockByNumber"
	getTransactionByHashMethod = "eth_getTransactionByHash"
)

type getblockClient struct {
	rpcURL string
	client *jsonrpc.Client
}

func New(c *jsonrpc.Client, token string) *getblockClient {
	return &getblockClient{
		rpcURL: fmt.Sprintf("%s/%s", getBlockBase, token),
		client: c,
	}
}

func (c *getblockClient) BlockNumber(ctx context.Context) (int64, error) {
	req := jsonrpc.NewRequest(blockNumberMethod, []any{}, getBlockID)

	resp, err := c.client.Call(ctx, c.rpcURL, req)
	if err != nil {
		return 0, err
	}

	blocknumberHex := resp.Result.(string)
	block, err := strconv.ParseInt(blocknumberHex, 0, 64)
	if err != nil {
		return 0, err
	}

	return block, nil
}

func (c *getblockClient) ChainID(ctx context.Context) (int64, error) {
	req := jsonrpc.NewRequest(chainIDMethod, []any{}, getBlockID)

	resp, err := c.client.Call(ctx, c.rpcURL, req)
	if err != nil {
		return 0, err
	}

	chainHex := resp.Result.(string)
	chain, err := strconv.ParseInt(chainHex, 0, 64)
	if err != nil {
		return 0, err
	}

	return chain, nil
}

func (c *getblockClient) BlockByNumber(ctx context.Context, number int64) (types.Block, error) {
	params := []any{
		fmt.Sprintf("0x%x", number),
		true,
	}

	req := jsonrpc.NewRequest(getBlockByNumberMethod, params, getBlockID)

	resp, err := c.client.Call(ctx, c.rpcURL, req)
	if err != nil {
		return types.Block{}, fmt.Errorf("call rpc %s: %w", getBlockByNumberMethod, err)
	}

	block, err := blockFromRPC(resp)
	if err != nil {
		return types.Block{}, fmt.Errorf("failed parse block from rpc: %w", err)
	}

	return block, nil
}

func (c *getblockClient) TransactionByHash(ctx context.Context, hash string) (types.Transaction, error) {
	params := []any{
		hash,
	}

	req := jsonrpc.NewRequest(getTransactionByHashMethod, params, getBlockID)

	resp, err := c.client.Call(ctx, c.rpcURL, req)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("%s: %w", getTransactionByHashMethod, err)
	}

	tx, err := txFromRPC(resp)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed parse tx from rpc: %w", err)
	}

	return tx, nil
}
