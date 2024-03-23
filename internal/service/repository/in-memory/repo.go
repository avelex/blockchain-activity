package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/avelex/blockchain-activity/internal/types"
)

type repository struct {
	mu *sync.RWMutex
	m  map[string]types.Block
}

func New() *repository {
	return &repository{
		mu: &sync.RWMutex{},
		m:  make(map[string]types.Block),
	}
}

func (r *repository) GetBlock(_ context.Context, number int64) (types.Block, error) {
	numberHex := fmt.Sprintf("0x%x", number)

	r.mu.RLock()
	defer r.mu.RUnlock()

	block, ok := r.m[numberHex]
	if !ok {
		return types.Block{}, errors.New("not found")
	}

	return block, nil
}

func (r *repository) PutBlock(_ context.Context, block types.Block) error {
	r.mu.Lock()
	r.m[block.Number] = block
	r.mu.Unlock()

	return nil
}
