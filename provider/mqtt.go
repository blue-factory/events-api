package provider

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/microapis/messages-api/channel"
	messagesiot "github.com/microapis/messages-iot-api"
	"github.com/nats-io/nats.go"
	"github.com/stoewer/go-strcase"
)

const (
	// MqttName the provider name
	MqttName = "mqtt"
	// MqttAPIKey the mqtt api key
	MqttAPIKey = "MqttApiKey"
)

// MqttProvider ...
type MqttProvider struct {
	Root channel.Provider
	Conn mqtt.Client
}

// NewMqtt ...
func NewMqtt() *MqttProvider {
	p := &MqttProvider{
		Root: channel.Provider{
			Name:   MqttName,
			Params: make(map[string]string),
		},
	}

	p.Root.Params[MqttAPIKey] = ""

	return p
}

// Approve ...
func (mp *MqttProvider) Approve(*messagesiot.Message) error {
	return nil
}

// Deliver ...
func (mp *MqttProvider) Deliver(m *messagesiot.Message) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	token := mp.Conn.Publish(m.Topic, byte(1), false, b)
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

	return nil
}

// LoadEnv ...
func (mp *MqttProvider) LoadEnv() error {
	env := strings.ToUpper(strcase.SnakeCase(MqttAPIKey))
	value := os.Getenv("PROVIDER_" + env)
	if value == "" {
		return errors.New("PROVIDER_" + env + " env value not defined")
	}

	mp.Root.Params[MqttAPIKey] = value

	return nil
}

// ToProvider ...
func (mp *MqttProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   mp.Root.Name,
		Params: mp.Root.Params,
	}
}

// NewMqttConnection ...
func NewMqttConnection(address, token string) (*nats.EncodedConn, error) {
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
