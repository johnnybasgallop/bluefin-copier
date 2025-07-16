package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type TradeEvent struct {
	Type      string  `json:"type"`
	Symbol    string  `json:"symbol"`
	Volume    float64 `json:"volume"`
	Price     float64 `json:"price"`
	Magic     int     `json:"magic"`
	Timestamp int64   `json:"timestamp"`  // original event time
	SentAt    int64   `json:"sent_at_ns"` // publish time
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// build the event
	evt := TradeEvent{
		Type:      "OPEN",
		Symbol:    "EURUSD",
		Volume:    1.0,
		Price:     1.08345,
		Magic:     12345,
		Timestamp: 1626420000,
		SentAt:    time.Now().UnixNano(), // record send time
	}

	// marshal to JSON
	b, err := json.Marshal(evt)
	if err != nil {
		panic(err)
	}

	// publish
	if err := rdb.Publish(ctx, "master:2001:events", b).Err(); err != nil {
		panic(fmt.Sprintf("failed to publish: %v", err))
	}

	fmt.Println("TradeEvent published")
}
