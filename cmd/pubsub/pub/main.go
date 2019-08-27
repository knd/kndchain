package main

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

func main() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	conn.Do("PUBLISH", "HELLO_CHANNEL", "message random")
}
