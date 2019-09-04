package mining

import (
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/wallet"
)

// Miner provides entry to mining actions
type Miner interface {
	Mine() error
}

type miner struct {
	service         mining.Service
	transactionPool wallet.TransactionPool
	wal             wallet.Wallet
	comm            pubsub.Service
}

// NewMiner creates a miner with necessary dependencies
func NewMiner(s mining.Service, p wallet.TransactionPool, w wallet.Wallet, c pubsub.Service) Miner {
	return &miner{s, p, w, c}
}

func (m *miner) Mine() error {
	// TODO: Implement this
	return nil
}
