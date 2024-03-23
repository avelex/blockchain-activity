package getblock_test

import (
	"context"
	"errors"
	"testing"

	"github.com/avelex/blockchain-activity/internal/adapters/rpc/getblock"
	"github.com/avelex/blockchain-activity/pkg/jsonrpc"
)

type mockRPCClient struct {
	call func(ctx context.Context, url string, request jsonrpc.Request) (jsonrpc.Response, error)
}

func (c *mockRPCClient) Call(ctx context.Context, url string, request jsonrpc.Request) (jsonrpc.Response, error) {
	return c.call(ctx, url, request)
}

func TestBlockNumber(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc            string
		rpc             getblock.JSONRPCClient
		wantBlockNumber int64
		err             error
	}{
		{
			desc: "Normal",
			rpc: &mockRPCClient{
				call: func(ctx context.Context, url string, request jsonrpc.Request) (jsonrpc.Response, error) {
					return jsonrpc.Response{
						Result: "0x1000",
					}, nil
				},
			},
			wantBlockNumber: 4096,
			err:             nil,
		},
		{
			desc: "RPC got error",
			rpc: &mockRPCClient{
				call: func(ctx context.Context, url string, request jsonrpc.Request) (jsonrpc.Response, error) {
					return jsonrpc.Response{}, errors.ErrUnsupported
				},
			},
			wantBlockNumber: 0,
			err:             errors.ErrUnsupported,
		},
		{
			desc: "Block number not string",
			rpc: &mockRPCClient{
				call: func(ctx context.Context, url string, request jsonrpc.Request) (jsonrpc.Response, error) {
					return jsonrpc.Response{
						Result: 4096,
					}, nil
				},
			},
			wantBlockNumber: 0,
			err:             getblock.ErrInvalidResult,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := getblock.New(tC.rpc, "")

			gotBlockNumber, err := c.BlockNumber(context.Background())
			if tC.err != nil && !errors.Is(err, tC.err) {
				t.FailNow()
			}

			if gotBlockNumber != tC.wantBlockNumber {
				t.Logf("Wrong block number, expected=%v got=%v", tC.wantBlockNumber, gotBlockNumber)
				t.FailNow()
			}
		})
	}
}
