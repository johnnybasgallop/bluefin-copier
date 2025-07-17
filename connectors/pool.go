package connectors

import (
	"fmt"
	"sync"

	"github.com/johnnybasgallop/bluefin-copier/internal/config"
)

// Pool manages a set of SlaveConnector instances, keyed by subscriber ID.
// It allows the Dispatcher to lookup the right connector for each CopyOrder.
type Pool struct {
	mu    sync.RWMutex               // protects access to conns map
	conns map[string]*SlaveConnector // subscriberID -> connector
}

// NewPool creates and initializes a Pool for all subscribers in the config.
// It instantiates one SlaveConnector per subscriber.
func NewPool(subscribers []config.Subscriber) (*Pool, error) {
	p := &Pool{
		conns: make(map[string]*SlaveConnector, len(subscribers)),
	}

	// Loop over each subscriber in your config
	for _, subCfg := range subscribers {
		// Create a connector for this subscriber
		conn, err := NewSlaveConnector(subCfg)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to initialize connector for subscriber %s: %w",
				subCfg.ID, err,
			)
		}
		p.conns[subCfg.ID] = conn
	}

	return p, nil
}

// Get retrieves the SlaveConnector for a given subscriber ID.
// Returns an error if no connector exists.
func (p *Pool) Get(subscriberID string) (*SlaveConnector, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	conn, ok := p.conns[subscriberID]
	if !ok {
		return nil, fmt.Errorf("no connector found for subscriber %s", subscriberID)
	}
	return conn, nil
}
