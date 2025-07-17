package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type TradeEvent struct {
	Type      string  `json:"type"`
	Symbol    string  `json:"symbol"`
	Volume    float64 `json:"volume"`
	Price     float64 `json:"price"`
	Magic     int     `json:"magic"`
	Timestamp int64   `json:"timestamp"`
	SentAt    int64   `json:"sent_at_ns"`
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	pubsub := rdb.Subscribe(ctx, "copier:commands")
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
		log.Fatalf("subscribe failed: %v", err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		// parse JSON
		var evt TradeEvent
		if err := json.Unmarshal([]byte(msg.Payload), &evt); err != nil {
			log.Printf("invalid JSON: %v", err)
			continue
		}

		// compute latency
		fmt.Printf(
			"Received %s %s\n",
			evt.Type,
			evt.Symbol,
		)
	}
}
