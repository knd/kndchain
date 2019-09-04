package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
	"github.com/knd/kndchain/pkg/wallet"
)

// Handler provides list of routes and action handlers
func Handler(l listing.Service, m mining.Service, c pubsub.Service, p wallet.TransactionPool, wal wallet.Wallet) http.Handler {
	router := httprouter.New()

	router.GET("/api/blocks", getBlocks(l))
	router.POST("/api/blocks", mineBlock(m, l, c))
	router.POST("/api/transactions", addTx(p, wal, c))
	router.GET("/api/transactions", getTxPool(p))

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

// type txInput struct {
// 	Timestamp int64  `json:"timestamp"`
// 	Amount    uint64 `json:"amount"`
// 	Address   string `json:"address"`
// 	Signature string `json:"sig"`
// }

// type txOutput struct {
// 	ID     string            `json:"id"`
// 	Output map[string]uint64 `json:"output"`
// 	Input  txInput           `json:"input"`
// }

func addTx(p wallet.TransactionPool, wal wallet.Wallet, c pubsub.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
			tx, err = wal.CreateTransaction(ati.Receiver, ati.Amount)
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

// func toTxOuptut(tx wallet.Transaction) txOutput {
// 	o := txOutput{}
// 	o.ID = tx.GetID()
// 	o.Output = tx.GetOutput()
// 	i := tx.GetInput()
// 	o.Input.Timestamp = i.Timestamp
// 	o.Input.Amount = i.Amount
// 	o.Input.Address = i.Address
// 	o.Input.Signature = hex.EncodeToString(i.Signature)

// 	return o
// }
