package pubsub

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
	"github.com/knd/kndchain/pkg/wallet"
)

const (
	// ChannelPubSub is channel name for publisher to send and subscribers to receive messages
	ChannelPubSub = "kndchain"

	// ChannelTransactions is channel name for publisher to send and subscribers to receive transactions
	ChannelTransactions = "kndtransactions"

	// PortPubSub is port on which pubsub server is run
	PortPubSub = ":6379"
)

// Service provides networking operations
type Service interface {
	Connect() error
	Disconnect() error
	SubscribePeers() error
	BroadcastBlockchain(bc *listing.Blockchain) error
	BroadcastTransaction(tx wallet.Transaction) error
}

type service struct {
	l   listing.Service
	m   mining.Service
	p   wallet.TransactionPool
	psc *redis.PubSubConn
}

// NewService creates a networking service with necessary dependencies
func NewService(l listing.Service, m mining.Service, p wallet.TransactionPool) Service {
	return &service{
		l: l,
		m: m,
		p: p,
	}
}

// Connect creates communication line with peers
func (s *service) Connect() error {
	conn, err := redis.Dial("tcp", PortPubSub)
	if err != nil {
		log.Printf("PubSubService#Connect: Failed to dial tcp connection, %v", err)
		return err
	}

	s.psc = &redis.PubSubConn{Conn: conn}

	return nil
}

// Disconnect closes communication line with peers
func (s *service) Disconnect() error {
	return s.psc.Conn.Close()
}

// BroadcastBlockchain broadcasts latest blockchain to peers
func (s *service) BroadcastBlockchain(bc *listing.Blockchain) error {
	b, err := json.Marshal(*bc)
	if err != nil {
		log.Printf("PubSubService#BroadcastBlockchain: Failed to json marshal blockchain, %v", err)
		return err
	}

	conn, err := redis.Dial("tcp", PortPubSub)
	defer conn.Close()
	if err != nil {
		log.Printf("PubSubService#BroadcastBlockchain: Failed to dial tcp connection, %v", err)
		return err
	}

	_, err = conn.Do("PUBLISH", ChannelPubSub, string(b[:]))

	return err
}

// BroadcastTransaction broadcasts latest transaction to peers
func (s *service) BroadcastTransaction(tx wallet.Transaction) error {
	b, err := json.Marshal(tx)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := redis.Dial("tcp", PortPubSub)
	defer conn.Close()
	if err != nil {
		log.Printf("PubSubService#BroadcastTransaction: Failed to dial tcp connection, %v", err)
		return err
	}

	_, err = conn.Do("PUBLISH", ChannelTransactions, string(b[:]))

	return err
}

// SubscribePeers listens to peers for incoming blockchain and transactions
func (s *service) SubscribePeers() error {
	err := s.psc.Subscribe(ChannelPubSub)
	if err != nil {
		log.Printf("PubSubService#SubscribePeers: Failed to subscribe to peers, channel=%s, %v", ChannelPubSub, err)
		return err
	}

	err = s.psc.Subscribe(ChannelTransactions)
	if err != nil {
		log.Printf("PubSubService#SubscribePeers: Failed to subscribe to peers, channel=%s, %v", ChannelTransactions, err)
		return err
	}

	go func(conn redis.Conn) {
		for conn.Err() == nil {
			switch v := s.psc.Receive().(type) {
			case redis.Message:
				if v.Channel == ChannelPubSub {
					// Received incoming blockchain
					// Replace incoming blockchain if valid longest chain
					var bc mining.Blockchain
					var err error
					err = json.Unmarshal(v.Data, &bc)
					if err != nil {
						log.Printf("PubSubService#SubscribePeers: Coudn't unmarshall incoming blockchain, %v", err)
						continue
					}

					err = s.m.ReplaceChain(&bc)
					if err != nil {
						log.Printf("PubSubService#SubscribePeers: Failed to replace blockchain, %v", err)
						continue
					}

					err = s.p.ClearBlockTransactions()
					if err != nil {
						log.Printf("PubSubService#SubscribePeers: Failed to clear block transactions in transaction pool, %v", err)
						continue
					}

					log.Printf("Replaced with longer chain. New len: %d", s.l.GetBlockCount())
				} else if v.Channel == ChannelTransactions {
					// Received incoming transaction
					// add transaction to pool
					var tx wallet.Tx
					var err error
					err = json.Unmarshal(v.Data, &tx)
					if err != nil {
						log.Printf("PubSubService#SubscribePeers: Couldn't unmarshall incoming transaction, %v", err)
						continue
					}
					s.p.Add(&tx)
					log.Printf("Transaction received. ID=%s", tx.GetID())
				}

			case redis.Subscription:
				log.Printf("Channel=%s, Kind=%s, Count=%d\n", v.Channel, v.Kind, v.Count)

			case error:
				log.Println(v)
			}
		}
	}(s.psc.Conn)

	return nil
}
