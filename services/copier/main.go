package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/johnnybasgallop/bluefin-copier/internal/config"
	"github.com/johnnybasgallop/bluefin-copier/internal/types"
)

func main() {
	// 1. Create a background context
	ctx := context.Background()

	// 2. Load your YAML config into a Config struct
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 3. Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 4. Subscribe to the master events channel
	channel := "master:" + "2001" + ":events"
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()

	// 5. Wait for the subscription to be active
	if _, err := sub.Receive(ctx); err != nil {
		log.Fatalf("subscribe failed: %v", err)
	}

	fmt.Printf("listening for events on %s …\n", channel)

	fmt.Printf("listening for events on %s …\n", channel)

	// Grab the Go channel that will deliver Redis messages
	ch := sub.Channel()

	// 6. Start processing incoming master events
	for {
		select {
		case <-ctx.Done():
			// If the context is cancelled (e.g. SIGINT), exit cleanly
			log.Println("copier shutting down")
			return

		case rawMsg, ok := <-ch:
			if !ok {
				// Channel closed, nothing more to read
				log.Println("subscription channel closed")
				return
			}

			// 6a. Parse the raw JSON into a TradeEvent object
			var evt types.TradeEvent
			if err := json.Unmarshal([]byte(rawMsg.Payload), &evt); err != nil {
				log.Printf("invalid TradeEvent JSON: %v", err)
				continue // skip to the next message
			}

			// 6b. For each subscriber defined in your config…
			for _, subCfg := range cfg.Subscribers {
				// i. Determine the output symbol (apply mapping if provided)
				outSym, ok := subCfg.SymbolMap[evt.Symbol]
				if !ok {
					outSym = evt.Symbol
				}

				// ii. Scale the volume by the subscriber’s lotRatio
				outVol := evt.Volume * subCfg.LotRatio

				// iii. (Optional) Apply filters, e.g. skip tiny volumes
				if outVol <= 0 {
					continue
				}

				// iv. Build a CopyOrder struct for this subscriber
				order := types.CopyOrder{
					SubscriberID:  subCfg.ID,
					Action:        evt.Type,
					Symbol:        outSym,
					Volume:        outVol,
					Price:         &evt.Price,                    // pointer for JSON omitempty support
					CorrelationID: fmt.Sprintf("%d", evt.SentAt), // track back to the original event
				}

				// v. Marshal the CopyOrder to JSON bytes
				data, err := json.Marshal(order)
				if err != nil {
					log.Printf("failed to marshal CopyOrder for %s: %v", subCfg.ID, err)
					continue
				}

				// vi. Publish the CopyOrder JSON to the dispatcher channel
				if err := rdb.Publish(ctx, "copier:commands", data).Err(); err != nil {
					log.Printf("failed to publish CopyOrder for %s: %v", subCfg.ID, err)
				}
			}
		}
	}

}
