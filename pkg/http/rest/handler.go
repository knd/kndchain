package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/knd/kndchain/pkg/crypto"

	"github.com/knd/kndchain/pkg/wallet"

	"github.com/julienschmidt/httprouter"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/networking/pubsub"
)

// Handler provides list of routes and action handlers
func Handler(l listing.Service, m mining.Service, c pubsub.Service, p wallet.TransactionPool) http.Handler {
	router := httprouter.New()

	router.GET("/api/blocks", getBlocks(l))
	router.POST("/api/blocks", mineBlock(m, l, c))
	router.POST("/api/transactions", addTx(p))

	return router
}

func getBlocks(l listing.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(l.GetBlockchain())
	}
}

// MineBlockInput encapsulates data/txs in new block
type MineBlockInput struct {
	Data []string `json:"data"`
}

func mineBlock(m mining.Service, l listing.Service, c pubsub.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var mineBlockInput MineBlockInput
		err := decoder.Decode(&mineBlockInput)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		lb := l.GetLastBlock()
		mb := &mining.Block{
			Timestamp:  lb.Timestamp,
			LastHash:   lb.LastHash,
			Hash:       lb.Hash,
			Data:       lb.Data,
			Nonce:      lb.Nonce,
			Difficulty: lb.Difficulty,
		}
		newBlock, err := m.MineNewBlock(mb, mineBlockInput.Data)
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

type addTxInput struct {
	Receiver string `json:"receiver"`
	Amount   uint64 `json:"amount"`
}

type txInput struct {
	Timestamp int64  `json:"timestamp"`
	Amount    uint64 `json:"amount"`
	Address   string `json:"address"`
	Signature string `json:"sig"`
}

type addTxOutput struct {
	ID     string            `json:"id"`
	Output map[string]uint64 `json:"output"`
	Input  txInput           `json:"input"`
}

func addTx(p wallet.TransactionPool) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)

		var txInput addTxInput
		err := decoder.Decode(&txInput)
		if err != nil || len(txInput.Receiver) == 0 || txInput.Amount <= 0 {
			http.Error(w, fmt.Sprintf("Invalid input err=%s, receiver=%s, amount=%d", err, txInput.Receiver, txInput.Amount), http.StatusBadRequest)
			return
		}

		// TODO: Replace wallet
		wallet := wallet.NewWallet(crypto.NewSecp256k1Generator())
		tx, err := wallet.CreateTransaction(txInput.Receiver, txInput.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		p.Add(tx)

		// marshalling
		o := addTxOutput{}
		o.ID = tx.GetID()
		o.Output = tx.GetOutput()
		i := tx.GetInput()
		o.Input.Timestamp = i.Timestamp
		o.Input.Amount = i.Amount
		o.Input.Address = i.Address
		o.Input.Signature = hex.EncodeToString(i.Signature)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(o)
	}
}
