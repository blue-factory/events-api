package provider

import (
	"errors"
	"os"
	"strings"

	"github.com/microapis/messages-lib/channel"
	messagesiot "github.com/microapis/messages-iot-api"
	"github.com/nats-io/nats.go"
	"github.com/stoewer/go-strcase"
)

const (
	// NatsName the provider name
	NatsName = "nats"
	// NatsAPIKey the nats api key
	NatsAPIKey = "NatsApiKey"
)

// NatsProvider ...
type NatsProvider struct {
	Root channel.Provider
	Conn *nats.EncodedConn
}

// NewNats ...
func NewNats() *NatsProvider {
	p := &NatsProvider{
		Root: channel.Provider{
			Name:   NatsName,
			Params: make(map[string]string),
		},
	}

	p.Root.Params[NatsAPIKey] = ""

	return p
}

// Approve ...
func (np *NatsProvider) Approve(*messagesiot.Message) error {
	return nil
}

// Deliver ...
func (np *NatsProvider) Deliver(m *messagesiot.Message) error {
	defer np.Conn.Close()

	err := np.Conn.Publish(m.Topic, m)
	if err != nil {
		return err
	}

	return nil
}

// LoadEnv ...
func (np *NatsProvider) LoadEnv() error {
	env := strings.ToUpper(strcase.SnakeCase(MqttAPIKey))
	value := os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	np.Root.Params[MqttAPIKey] = value

	return nil
}

// ToProvider ...
func (np *NatsProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   np.Root.Name,
		Params: np.Root.Params,
	}
}

// NewNatsConnection ...
func NewNatsConnection(address, token string) (*nats.EncodedConn, error) {
	natsConn, err := nats.Connect(address, nats.Token(token))
	if err != nil {
		return nil, err
	}
	c, err := nats.NewEncodedConn(natsConn, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return c, nil
}
