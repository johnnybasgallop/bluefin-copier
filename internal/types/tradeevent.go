package types

// Represents the JSON you publish from your publisher
type TradeEvent struct {
	Type      string  `json:"type"`       // "OPEN", "CLOSE", etc.
	Symbol    string  `json:"symbol"`     // "EURUSD"
	Volume    float64 `json:"volume"`     // 1.0
	Price     float64 `json:"price"`      // 1.08345
	Magic     int     `json:"magic"`      // user-defined tag
	Timestamp int64   `json:"timestamp"`  // original event epoch
	SentAt    int64   `json:"sent_at_ns"` // nano-ts when published
}
