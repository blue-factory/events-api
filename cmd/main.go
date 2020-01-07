package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	mc "github.com/microapis/clients-go/messages"
	"github.com/microapis/messages-api/backend"
	"github.com/microapis/messages-api/channel"
	messagesiot "github.com/microapis/messages-iot-api"
	backendIot "github.com/microapis/messages-iot-api/backend"
	"github.com/microapis/messages-iot-api/provider"
)

var (
	nats *provider.NatsProvider
	mqtt *provider.MqttProvider
)

func main() {
	var err error

	// read providers env values
	providersEnv := os.Getenv("PROVIDERS")
	if providersEnv == "" {
		log.Fatal(errors.New("PROVIDERS value not defined"))
	}

	// define provider slice names
	ppn := strings.Split(providersEnv, ",")
	// define providers slice
	pp := make([]*channel.Provider, 0)

	// iterate over providers name
	for _, v := range ppn {
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
		log.Fatal(err)
	}

	// read messages-api env values
	messagesHost := os.Getenv("MESSAGES_HOST")
	if messagesHost == "" {
		log.Fatal(errors.New("MESSAGES_HOST value not defined"))
	}
	messagesPort := os.Getenv("MESSAGES_PORT")
	if messagesPort == "" {
		log.Fatal(errors.New("MESSAGES_PORT value not defined"))
	}

	// register channel on messages-api
	addr := fmt.Sprintf("%s:%s", messagesHost, messagesPort)

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal(errors.New("HOST value not defined"))
	}

	// get grpc port env value:
	port := os.Getenv("PORT")
	if port == "" {
		err := errors.New("invalid PORT env value")
		log.Fatal(err)
	}

	// create channel to register
	c := &channel.Channel{
		Name:      messagesiot.Channel,
		Host:      host,
		Port:      port,
		Providers: pp,
	}

	MessagesAPI, err := mc.NewMessagesAPI(addr, "token")
	if err != nil {
		log.Fatalln(err)
	}

	err = MessagesAPI.Register(c)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Channel %v is registered on messages-api with providers: %v", messagesiot.Channel, c.ProvidersNames())

	// define address value to grpc service
	addr = fmt.Sprintf("0.0.0.0:%s", port)

	// define service with Approve and Deliver methods
	svc := backendIot.NewBackend(nats, mqtt)

	// start grpc pigeon-ses-api service
	log.Printf("Serving at %s", addr)
	if err := backend.ListenAndServe(addr, svc); err != nil {
		log.Fatal(err)
	}
}
