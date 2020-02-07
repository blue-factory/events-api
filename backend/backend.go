package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/microapis/events-api"
	"github.com/microapis/events-api/provider"
	"github.com/microapis/events-api/provider/mqttprovider"
	"github.com/microapis/events-api/provider/natsprovider"
)

const (
	prefix = "PROVIDER_"
)

var (
	nats provider.Provider
	mqtt provider.Provider
)

// Backend ...
type Backend struct {
	Nats provider.Provider
	Mqtt provider.Provider
}

// NewBackend ...
func NewBackend(providers []string) (*Backend, error) {
	var err error

	// define repetitive msg error sufix
	errSufix := " env value not defined"

	// iterate over providers name
	for _, v := range providers {
		switch v {
		case natsprovider.Name:
			// get NATS_HOST env value
			env := "NATS_HOST"
			natsHost := os.Getenv(prefix + env)
			if natsHost == "" {
				return nil, errors.New(prefix + env + errSufix)
			}

			// get NATS_PORT env value
			env = "NATS_PORT"
			natsPort := os.Getenv(prefix + env)
			if natsPort == "" {
				return nil, errors.New(prefix + env + errSufix)
			}

			// get NATS_API_KEY env value
			env = "NATS_API_KEY"
			natsAPIKey := os.Getenv(prefix + env)
			if natsAPIKey == "" {
				return nil, errors.New(prefix + env + errSufix)
			}

			address := fmt.Sprintf("%s:%s", natsHost, natsPort)
			nats, err = natsprovider.New(address, natsAPIKey)
			if err != nil {
				return nil, err
			}
		case mqttprovider.Name:
			// get MQTT_HOST env value
			env := "MQTT_HOST"
			mqttHost := os.Getenv(prefix + env)
			if mqttHost == "" {
				return nil, errors.New(prefix + env + errSufix)
			}

			// get MQTT_PORT env value
			env = "MQTT_PORT"
			mqttPort := os.Getenv(prefix + env)
			if mqttPort == "" {
				return nil, errors.New(prefix + env + errSufix)
			}

			// get MQTT_API_KEY env value
			env = "MQTT_API_KEY"
			mqttAPIKey := os.Getenv(prefix + env)
			if mqttAPIKey == "" {
				return nil, errors.New(prefix + env + errSufix)
			}

			address := fmt.Sprintf("%s:%s", mqttHost, mqttPort)
			mqtt, err = mqttprovider.New(address, mqttAPIKey)
			if err != nil {
				return nil, err
			}
		}
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

	m := new(events.Message)
	err = json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Error] error = %v", err))
		return false, err
	}

	log.Println(1111)
	switch m.Provider {
	case natsprovider.Name:
		err = b.Nats.Approve(m)
	case mqttprovider.Name:
		err = b.Mqtt.Approve(m)
	default:
		return false, errors.New("invalid provider value")
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Error] error = %v", err))
		return false, err
	}

	log.Println(222)

	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Approve][Response] message = %v", m))

	return true, nil
}

// Deliver ...
func (b *Backend) Deliver(content string) error {
	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Request] content = %v", content))
	log.Println(111)
	if content == "" {
		err := errors.New("Invalid message content")
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Error] error = %v", err))
		return err
	}
	log.Println(222)

	m := new(events.Message)
	err := json.Unmarshal([]byte(content), m)
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Error] error = %v", err))
		return err
	}
	log.Println(333)

	switch m.Provider {
	case natsprovider.Name:
		err = b.Nats.Deliver(m)
	case mqttprovider.Name:
		err = b.Mqtt.Deliver(m)
	}
	if err != nil {
		log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Error] error = %v", err))
		return err
	}
	log.Println(444)

	log.Println(fmt.Sprintf("[gRPC][MessagesIoTService][Deliver][Response] message = %v", m))

	return nil
}
