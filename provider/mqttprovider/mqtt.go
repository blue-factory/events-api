package mqttprovider

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/microapis/messages-core/channel"
	"github.com/microapis/events-api"
)

const (
	// Name the provider name
	Name = "mqtt"
)

// MqttProvider ...
type MqttProvider struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
	Conn   mqtt.Client
}

// New ...
func New(address string, token string) (*MqttProvider, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", address))

	// opts.SetUsername(uri.User.Username())
	// password, _ := uri.User.Password()
	// opts.SetPassword(password)

	opts.SetClientID("client-id")
	client := mqtt.NewClient(opts)
	t := client.Connect()
	for !t.WaitTimeout(3 * time.Second) {
	}
	if err := t.Error(); err != nil {
		return nil, err
	}

	p := &MqttProvider{
		Name:   Name,
		Params: make(map[string]string),
		Conn:   client,
	}

	return p, nil
}

// Approve ...
func (mp *MqttProvider) Approve(*events.Message) error {
	return nil
}

// Deliver ...
func (mp *MqttProvider) Deliver(m *events.Message) error {
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

// ToProvider ...
func (mp *MqttProvider) ToProvider() *channel.Provider {
	return &channel.Provider{
		Name:   mp.Name,
		Params: mp.Params,
	}
}
