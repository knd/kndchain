package pubsub

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/knd/kndchain/pkg/listing"
	"github.com/knd/kndchain/pkg/mining"
)

const (
	// ChannelPubSub is channel name for publisher to send and subscribers to receive messages
	ChannelPubSub = "kndchain"

	// PortPubSub
	PortPubSub = ":6379"
)

// Service provides networking operations
type Service interface {
	Connect() error
	Disconnect() error
	SubscribePeers() error
	BroadcastBlockchain(bc *listing.Blockchain) error
}

type service struct {
	l   listing.Service
	m   mining.Service
	psc *redis.PubSubConn
}

// NewService creates a networking service with necessary dependencies
func NewService(l listing.Service, m mining.Service) Service {
	return &service{
		l: l,
		m: m,
	}
}

// Connect creates communication line with peers
func (s *service) Connect() error {
	conn, err := redis.Dial("tcp", PortPubSub)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	conn, err := redis.Dial("tcp", PortPubSub)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	_, err = conn.Do("PUBLISH", ChannelPubSub, string(b[:]))

	return err
}

// SubscribePeers listens to peers for incoming blockchain
func (s *service) SubscribePeers() error {
	err := s.psc.Subscribe(ChannelPubSub)
	if err != nil {
		log.Fatal(err)
	}

	go func(conn redis.Conn) {
		for conn.Err() == nil {
			switch v := s.psc.Receive().(type) {
			case redis.Message:
				log.Printf("%s: message: %s\n", v.Channel, v.Data)
			case redis.Subscription:
				log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				log.Println(v)
			}
		}
	}(s.psc.Conn)

	return nil
}
