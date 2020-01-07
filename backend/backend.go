package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/microapis/messages-iot-api"
	"github.com/microapis/messages-iot-api/provider"
)

// Backend ...
type Backend struct {
	Nats *provider.NatsProvider
	Mqtt *provider.MqttProvider
}

// NewBackend ...
func NewBackend(nats *provider.NatsProvider, mqtt *provider.MqttProvider) *Backend {
	return &Backend{
		Nats: nats,
		Mqtt: mqtt,
	}
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
