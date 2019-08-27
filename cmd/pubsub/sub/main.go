package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

func main() {
	for {
		conn, err := redis.Dial("tcp", ":6379")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		psc := redis.PubSubConn{Conn: conn}
		psc.Subscribe("HELLO_CHANNEL")

		for conn.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.Message:
				fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Println(v)
			}
		}
	}
}
