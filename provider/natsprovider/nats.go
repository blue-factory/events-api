package natsprovider

import (
	"log"
	"time"

	"github.com/microapis/events-api"
	"github.com/microapis/messages-core/channel"
	"github.com/nats-io/nats.go"
)

const (
	// Name the provider name
	Name = "nats"
)

// NatsProvider ...
type NatsProvider struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
	Conn   *nats.EncodedConn
}

// New ...
func New(address string, token string) (*NatsProvider, error) {
	natsConn, err := nats.Connect("nats://"+address, nats.MaxReconnects(15), nats.ReconnectWait(3*time.Second))
	if err != nil {
		return nil, err
	}
	c, err := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	p := &NatsProvider{
		Name:   Name,
		Params: make(map[string]string),
		Conn:   c,
	}

	return p, nil
}

// Approve ...
func (np *NatsProvider) Approve(*events.Message) error {
	return nil
}

// Deliver ...
func (np *NatsProvider) Deliver(m *events.Message) error {
	defer np.Conn.Close()

	log.Println("====")
	log.Println(m.Topic, m)
	log.Println("====")

	err := np.Conn.Publish(m.Topic, m)
	if err != nil {
		return err
	}

	defer np.Conn.Close()
	return nil
}

// ToProvider ...
func (np *NatsProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   np.Name,
		Params: np.Params,
	}
}
