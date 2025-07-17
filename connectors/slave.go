// connectors/slave.go
package connectors

import (
	"context"
	"fmt"

	"github.com/johnnybasgallop/bluefin-copier/internal/config"
	"github.com/johnnybasgallop/bluefin-copier/internal/types"
)

// SlaveConnector holds the broker session for one subscriber account.
// It knows how to send orders to that subscriber's broker.
type SlaveConnector struct {
	SubscriberID string        // identifier for the subscriber account
	session      BrokerSession // abstract broker protocol client
}

// BrokerSession defines the subset of methods we need from a broker client.
// In real code, replace these placeholder request/response types
// with your actual protocol types (e.g. FIX NewOrderSingle).
type BrokerSession interface {
	// SendNewOrder sends an order request and returns the broker's response.
	SendNewOrder(ctx context.Context, req *NewOrderRequest) (*NewOrderResponse, error)
}

// NewOrderRequest represents the data needed to open a new order on the broker.
// Adapt this to your broker's API or protocol.
type NewOrderRequest struct {
	AccountID     string   // which subscriber account to trade on
	Symbol        string   // e.g. "EURUSD" or mapped symbol
	Side          string   // typically "BUY" or "SELL"
	OrderQty      float64  // size of the order
	Price         *float64 // limit price (nil for market)
	ClientOrderID string   // unique ID for idempotency/tracing
}

// NewOrderResponse captures the broker's reply to a new-order request.
type NewOrderResponse struct {
	Status       string // e.g. "Filled", "Rejected"
	RejectReason string // non-empty if Status == "Rejected"
}

// NewSlaveConnector creates and initializes a SlaveConnector for a single subscriber.
// It should authenticate and establish the broker session.
func NewSlaveConnector(subCfg config.Subscriber) (*SlaveConnector, error) {
	// TODO: Replace the fake session with real broker session creation:
	// e.g. session, err := brokerproto.NewSession(subCfg.Credentials...)
	session, err := NewFakeSession(subCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect broker for subscriber %s: %w", subCfg.ID, err)
	}
	return &SlaveConnector{
		SubscriberID: subCfg.ID,
		session:      session,
	}, nil
}

// SendOrder sends the given CopyOrder to the broker via the Session.
// It translates the CopyOrder into a NewOrderRequest and checks the response.
func (c *SlaveConnector) SendOrder(ctx context.Context, order types.CopyOrder) error {
	// Build the broker request
	req := &NewOrderRequest{
		AccountID:     c.SubscriberID,
		Symbol:        order.Symbol,
		Side:          order.Action,
		OrderQty:      order.Volume,
		Price:         order.Price,
		ClientOrderID: order.CorrelationID,
	}

	// Send it over the broker session
	resp, err := c.session.SendNewOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("error sending order for %s: %w", c.SubscriberID, err)
	}

	// Check for broker-level rejection
	if resp.Status == "Rejected" {
		return fmt.Errorf("broker rejected order for %s: %s", c.SubscriberID, resp.RejectReason)
	}

	// On success, resp.Status might be "Filled" or "Accepted"
	return nil
}

// --- Fake session for local testing purposes only ---
// Remove this when wiring up a real broker client.

// NewFakeSession returns a dummy session that always succeeds.
func NewFakeSession(_ config.Subscriber) (BrokerSession, error) {
	return &fakeSession{}, nil
}

type fakeSession struct{}

func (s *fakeSession) SendNewOrder(ctx context.Context, req *NewOrderRequest) (*NewOrderResponse, error) {
	// Simulate network delay or processing if you like:
	// time.Sleep(10 * time.Millisecond)
	return &NewOrderResponse{Status: "Filled"}, nil
}
