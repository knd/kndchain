package calculating

import (
	"log"

	"github.com/knd/kndchain/pkg/config"
)

// Service provides access to calculating operations
type Service interface {
	Balance(address string, bc *Blockchain) uint64
	BalanceByBlockIndex(address string, bc *Blockchain, i int) uint64
}

type service struct{}

// NewService creates a calculating service
func NewService() Service {
	return &service{}
}

// Balance returns the current balance of the address given blockchain history
func (s *service) Balance(address string, bc *Blockchain) uint64 {
	return s.BalanceByBlockIndex(address, bc, len(bc.Chain)-1)
}

func (s *service) BalanceByBlockIndex(address string, bc *Blockchain, index int) uint64 {
	var balance uint64
	if bc == nil || len(bc.Chain) == 0 {
		return config.InitialBalance
	}
	if index >= len(bc.Chain) {
		log.Fatalf("Index is greater than block count index=%d", index)
	}
	if index < 0 {
		log.Fatalf("Index is less than zero index=%d", index)
	}

	var foundWalletTxInBlock bool
	for i := index; i >= 0; i-- {
		block := bc.Chain[i]

		var blockAmount uint64
		for _, tx := range block.Data {
			if tx.Input.Address == address {
				blockAmount = tx.Output[address]
				foundWalletTxInBlock = true
			} else if amount, ok := tx.Output[address]; ok {
				blockAmount += amount
			}
		}

		balance += blockAmount

		if foundWalletTxInBlock {
			break
		}
	}

	if foundWalletTxInBlock {
		return balance
	}

	return config.InitialBalance + balance
}
