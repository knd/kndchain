package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/knd/kndchain/pkg/calculating"

	"github.com/julienschmidt/httprouter"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/miner"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/wallet"
)

// Handler provides list of routes and action handlers
func Handler(l listing.Service, m mining.Service, c pubsub.Service, p wallet.TransactionPool, wal wallet.Wallet, miner miner.Miner, cal calculating.Service, enableMiningAPI bool) http.Handler {
	router := httprouter.New()

	router.GET("/api/blocks", getBlocks(l))
	router.POST("/api/blocks", mineBlock(m, l, c))
	router.POST("/api/transactions", addTx(p, wal, c, l))
	router.GET("/api/transactions", getTxPool(p))
	if enableMiningAPI {
		router.GET("/api/mine-transactions", mineTransactions(miner, l))     // for testing
		router.GET("/api/mining-address", getMiningAddressInfo(wal, l, cal)) // for testing
	}

	return router
}

func getBlocks(l listing.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(l.GetBlockchain())
	}
}

func mineBlock(m mining.Service, l listing.Service, c pubsub.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var transaction mining.Transaction
		err := decoder.Decode(&transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lb := l.GetLastBlock()
		mb := &mining.Block{
			Timestamp:  lb.Timestamp,
			LastHash:   lb.LastHash,
			Hash:       lb.Hash,
			Data:       toMiningTransactions(lb.Data),
			Nonce:      lb.Nonce,
			Difficulty: lb.Difficulty,
		}
		newBlock, err := m.MineNewBlock(mb, []mining.Transaction{transaction})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = m.AddBlock(newBlock)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c.BroadcastBlockchain(l.GetBlockchain())

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(l.GetLastBlock())
	}
}

func toMiningTransactions(data []listing.Transaction) []mining.Transaction {
	var mTxs []mining.Transaction
	for _, transaction := range data {
		mTxs = append(mTxs, mining.Transaction{
			ID:     transaction.ID,
			Output: transaction.Output,
			Input: mining.Input{
				Timestamp: transaction.Input.Timestamp,
				Amount:    transaction.Input.Amount,
				Address:   transaction.Input.Address,
				Signature: transaction.Input.Signature,
			},
		})
	}
	return mTxs
}

type addTxInput struct {
	Receiver string `json:"receiver"`
	Amount   uint64 `json:"amount"`
}

func addTx(p wallet.TransactionPool, wal wallet.Wallet, c pubsub.Service, lister listing.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var ati addTxInput
		err := decoder.Decode(&ati)
		if err != nil || len(ati.Receiver) == 0 || ati.Amount <= 0 {
			http.Error(w, fmt.Sprintf("Invalid input err=%s, receiver=%s, amount=%d", err, ati.Receiver, ati.Amount), http.StatusBadRequest)
			return
		}

		var tx wallet.Transaction
		if p.Exists(wal.PubKeyHex()) {
			tx = p.GetTransaction(wal.PubKeyHex())
			err = tx.Append(wal, ati.Receiver, ati.Amount)
		} else {
			tx, err = wal.CreateTransaction(ati.Receiver, ati.Amount, lister)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.Add(tx)
		c.BroadcastTransaction(tx)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tx)
	}
}

func getTxPool(p wallet.TransactionPool) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		output := make(map[string]wallet.Transaction)

		for _, tx := range p.All() {
			output[tx.GetID()] = tx
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	}
}

func mineTransactions(miner miner.Miner, lister listing.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := miner.Mine()
		if err != nil {
			log.Printf("Failed to mine, %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lister.GetLastBlock())
	}
}

// AddressInfo is return result from getMiningAddressInfo
type AddressInfo struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
}

func getMiningAddressInfo(wal wallet.Wallet, lister listing.Service, cal calculating.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		addressInfo := AddressInfo{
			Address: wal.PubKeyHex(),
			Balance: cal.Balance(wal.PubKeyHex(), toCalculatingBlockchain(lister.GetBlockchain())),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(addressInfo)
	}
}

func toCalculatingBlockchain(bc *listing.Blockchain) *calculating.Blockchain {
	if bc == nil {
		return nil
	}

	result := &calculating.Blockchain{}
	for _, block := range bc.Chain {
		cTransactions := []calculating.Transaction{}
		for _, transaction := range block.Data {
			cTx := calculating.Transaction{
				ID:     transaction.ID,
				Output: transaction.Output,
				Input: calculating.Input{
					Timestamp: transaction.Input.Timestamp,
					Amount:    transaction.Input.Amount,
					Address:   transaction.Input.Address,
					Signature: transaction.Input.Signature,
				},
			}
			cTransactions = append(cTransactions, cTx)
		}
		cBlock := calculating.Block{
			Timestamp:  block.Timestamp,
			LastHash:   block.LastHash,
			Hash:       block.Hash,
			Data:       cTransactions,
			Nonce:      block.Nonce,
			Difficulty: block.Difficulty,
		}
		result.Chain = append(result.Chain, cBlock)
	}

	return result
}
