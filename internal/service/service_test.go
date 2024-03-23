package service_test

import (
	"context"
	"errors"
	"slices"
	"sort"
	"testing"

	"github.com/avelex/blockchain-activity/internal/service"
	inmemory "github.com/avelex/blockchain-activity/internal/service/repository/in-memory"
	"github.com/avelex/blockchain-activity/internal/types"
)

type mockBlockchainProvider struct {
	blockNumber       func(ctx context.Context) (int64, error)
	chainID           func(ctx context.Context) (int64, error)
	blockByNumber     func(ctx context.Context, number int64) (types.Block, error)
	transactionByHash func(ctx context.Context, hash string) (types.Transaction, error)
}

func (p *mockBlockchainProvider) BlockNumber(ctx context.Context) (int64, error) {
	return p.blockNumber(ctx)
}

func (p *mockBlockchainProvider) ChainID(ctx context.Context) (int64, error) {
	return p.chainID(ctx)
}

func (p *mockBlockchainProvider) BlockByNumber(ctx context.Context, number int64) (types.Block, error) {
	return p.blockByNumber(ctx, number)
}

func (p *mockBlockchainProvider) TransactionByHash(ctx context.Context, hash string) (types.Transaction, error) {
	return p.transactionByHash(ctx, hash)
}

func TestTopAddressesByActivity(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc     string
		provider service.BlockchainProvider
		cache    service.Cache
		want     types.BlockchainActivity
		err      error
	}{
		{
			desc: "Normal",
			provider: &mockBlockchainProvider{
				blockNumber: func(ctx context.Context) (int64, error) {
					return 4096, nil
				},
				chainID: func(ctx context.Context) (int64, error) {
					return 1, nil
				},
				blockByNumber: func(ctx context.Context, number int64) (types.Block, error) {
					return types.Block{
						Number: "0x1000",
						Transactions: []types.Transaction{
							{
								From:  "0x0000000000000000000000000000000000000001",
								To:    "0x0000000000000000000000000000000000000002",
								Input: "0xa9059cbb000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000bea5237b",
							},
							{
								From:  "0x0000000000000000000000000000000000000001",
								To:    "0x0000000000000000000000000000000000000003",
								Input: "0xa9059cbb000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000bea5237b",
							},
						},
					}, nil
				},
			},
			cache: inmemory.New(),
			want: types.BlockchainActivity{
				BlockNumber: 4096,
				Chain:       1,
				Addresses: []types.AddressActivity{
					{Address: "0x0000000000000000000000000000000000000001", Activity: 200},
					{Address: "0x0000000000000000000000000000000000000002", Activity: 100},
					{Address: "0x0000000000000000000000000000000000000003", Activity: 100},
				},
			},
			err: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			svc := service.New(tC.cache, tC.provider)

			got, err := svc.TopAddressesByActivity(context.Background())
			if tC.err != nil && !errors.Is(err, tC.err) {
				t.FailNow()
			}

			if tC.want.BlockNumber != got.BlockNumber {
				t.Logf("Wrong block number, expected=%v got=%v", tC.want.BlockNumber, got.BlockNumber)
				t.FailNow()
			}

			if tC.want.Chain != got.Chain {
				t.Logf("Wrong chain id, expected=%v got=%v", tC.want.Chain, got.Chain)
				t.FailNow()
			}

			sort.Slice(got.Addresses, func(i, j int) bool {
				return got.Addresses[i].Address < got.Addresses[j].Address
			})

			if !slices.EqualFunc(tC.want.Addresses, got.Addresses, func(aa1, aa2 types.AddressActivity) bool {
				if aa1.Address != aa2.Address {
					return false
				}

				if aa1.Activity != aa2.Activity {
					return false
				}

				return true
			}) {
				t.Logf("Wrong activity, expected=%v got=%v", tC.want.Addresses, got.Addresses)
				t.FailNow()
			}
		})
	}
}
