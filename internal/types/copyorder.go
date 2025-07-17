package types

// What the Copier emits for each subscriber
type CopyOrder struct {
	SubscriberID  string   `json:"subscriber_id"`
	Action        string   `json:"action"`          // from TradeEvent.Type
	Symbol        string   `json:"symbol"`          // mapped symbol
	Volume        float64  `json:"volume"`          // scaled volume
	Price         *float64 `json:"price,omitempty"` // nil for market orders
	CorrelationID string   `json:"correlation_id"`  // e.g. the SentAt or a UUID
}
