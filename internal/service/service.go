package service

import (
	"context"
	"log/slog"
	"slices"

	"github.com/avelex/blockchain-activity/internal/types"
	"github.com/avelex/blockchain-activity/internal/utils/erc20"
)

type BlockchainProvider interface {
	BlockNumber(ctx context.Context) (int64, error)
	ChainID(ctx context.Context) (int64, error)
	BlockByNumber(ctx context.Context, number int64) (types.Block, error)
	TransactionByHash(ctx context.Context, hash string) (types.Transaction, error)
}

type Cache interface {
	GetBlock(ctx context.Context, number int64) (types.Block, error)
	PutBlock(ctx context.Context, block types.Block) error
}

type blockchainService struct {
	cache    Cache
	provider BlockchainProvider
}

func New(cache Cache, provider BlockchainProvider) *blockchainService {
	return &blockchainService{
		cache:    cache,
		provider: provider,
	}
}

func (svc *blockchainService) TopAddressesByActivity(ctx context.Context) (types.BlockchainActivity, error) {
	n, err := svc.provider.BlockNumber(ctx)
	if err != nil {
		return types.BlockchainActivity{}, err
	}

	blocks, err := svc.getRangeBlocks(ctx, n, 100)
	if err != nil {
		return types.BlockchainActivity{}, err
	}

	addresses := topActiveAddressFromBlocks(blocks, 5)

	chainID, err := svc.provider.ChainID(ctx)
	if err != nil {
		return types.BlockchainActivity{}, err
	}

	return types.BlockchainActivity{
		BlockNumber: n,
		Chain:       chainID,
		Addresses:   addresses,
	}, nil
}

// fetch blocks information using workers pool pattern to increase speed
func (svc *blockchainService) getRangeBlocks(ctx context.Context, start, limit int64) ([]types.Block, error) {
	end := start - limit
	blocks := make([]types.Block, 0, limit)

	// TODO: move workers count to config
	workers := make(chan int64, 50)
	result := make(chan types.Block, 10)

	for i := 0; i < cap(workers); i++ {
		go func() {
			for blockNumber := range workers {
				if block, err := svc.cache.GetBlock(ctx, blockNumber); err == nil {
					result <- block
					continue
				}

				block, err := svc.provider.BlockByNumber(ctx, blockNumber)
				if err == nil {
					if err := svc.cache.PutBlock(ctx, block); err != nil {
						slog.Warn(err.Error())
					}
				}

				result <- block
			}
		}()
	}

	go func() {
		for i := start; i > end; i-- {
			workers <- i
		}
	}()

	for i := start; i > end; i-- {
		block := <-result
		if block.Number == "" {
			continue
		}

		blocks = append(blocks, block)
	}

	close(workers)
	close(result)

	slog.Info("Successful fetch", "block", len(blocks))

	return blocks, nil
}

func topActiveAddressFromBlocks(blocks []types.Block, topCount int) []types.AddressActivity {
	activity := make(map[string]uint)

	if topCount < 0 {
		return []types.AddressActivity{}
	}

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			if !erc20.IsTransfer(tx.Input) {
				continue
			}

			transfer, err := erc20.TransferFromTransactionInput(tx.Input, tx.From)
			if err != nil {
				slog.Warn(err.Error())
				continue
			}

			activity[transfer.Sender]++
			activity[transfer.Receiver]++
		}
	}

	keys := make([]string, 0, len(activity))
	for k := range activity {
		keys = append(keys, k)
	}

	slices.SortStableFunc(keys, func(a, b string) int {
		tempA := activity[a]
		tempB := activity[b]

		if tempA == tempB {
			return 0
		} else if tempA > tempB {
			return -1
		} else {
			return 1
		}
	})

	top := make([]types.AddressActivity, 0, topCount)

	topEnd := topCount
	if len(keys) < topCount {
		topEnd = len(keys)
	}

	for _, v := range keys[:topEnd] {
		top = append(top, types.AddressActivity{
			Address:  v,
			Activity: activity[v],
		})
	}

	return top
}
