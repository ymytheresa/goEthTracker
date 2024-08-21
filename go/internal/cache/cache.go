package cache

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Transfer struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Time   time.Time
}

type Cache struct {
	transfers []Transfer
	mu        sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		transfers: make([]Transfer, 0),
	}
}

func (c *Cache) AddTransfer(transfer Transfer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.transfers = append(c.transfers, transfer)
}

func (c *Cache) GetTransfers() []Transfer {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Transfer{}, c.transfers...)
}

func (c *Cache) GetTransfersForAddress(address common.Address) []Transfer {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []Transfer
	for _, t := range c.transfers {
		if t.From == address || t.To == address {
			result = append(result, t)
		}
	}
	return result
}
