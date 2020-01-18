package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	messagesiot "github.com/microapis/messages-iot-api"
	"github.com/microapis/messages-iot-api/provider"
	"github.com/microapis/messages-core/channel"
)

var (
	nats *provider.NatsProvider
	mqtt *provider.MqttProvider
)

// Backend ...
type Backend struct {
	Nats *provider.NatsProvider
	Mqtt *provider.MqttProvider
}

// NewBackend ...
func NewBackend(providers []string) (*Backend, error) {
	var err error

	// define providers slice
	pp := make([]*channel.Provider, 0)

	// iterate over providers name
	for _, v := range providers {
		switch v {
		case provider.NatsName:
			nats = provider.NewNats()
			err = nats.LoadEnv()
			pp = append(pp, nats.ToProvider())
		case provider.MqttName:
			mqtt = provider.NewMqtt()
			err = mqtt.LoadEnv()
			pp = append(pp, mqtt.ToProvider())
		}
	}
	if err != nil {
		return nil, err
	}

	return &Backend{
		Nats: nats,
		Mqtt: mqtt,
	}, nil
}

// Approve ...
func (b *Backend) Approve(content string) (valid bool, err error) {
	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Request] content = %v", content))

	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Error] error = %v", err))
		return false, err
	}

	m := new(messagesiot.Message)
	err = json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Error] error = %v", err))
		return false, err
	}

	switch m.Provider {
	case provider.NatsName:
		err = b.Nats.Approve(m)
	case provider.MqttName:
		err = b.Mqtt.Approve(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Error] error = %v", err))
		return false, err
	}

	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Response] message = %v", m))

	return true, nil
}

// Deliver ...
func (b *Backend) Deliver(content string) error {
	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Request] content = %v", content))

	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Error] error = %v", err))
		return err
	}

	m := new(messagesiot.Message)
	err := json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Error] error = %v", err))
		return err
	}

	switch m.Provider {
	case provider.NatsName:
		err = b.Nats.Deliver(m)
	case provider.MqttName:
		err = b.Mqtt.Deliver(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Error] error = %v", err))
		return err
	}

	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Response] message = %v", m))

	return nil
}
