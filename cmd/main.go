package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	backendIot "github.com/microapis/messages-iot-api/backend"
	"github.com/microapis/messages-core/service"
)

func main() {
	var err error

	// get providers env values
	providersEnv := os.Getenv("PROVIDERS")
	if providersEnv == "" {
		log.Fatal(errors.New("PROVIDERS value not defined"))
	}
	// define provider slice names
	providers := strings.Split(providersEnv, ",")

	// get grpc port env value
	port := os.Getenv("PORT")
	if port == "" {
		err := errors.New("invalid PORT env value")
		log.Fatal(err)
	}

	// get redis host env value
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		err := errors.New("invalid REDIS_HOST env value")
		log.Fatal(err)
	}

	// get redis port env value
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		err := errors.New("invalid REDIS_PORT env value")
		log.Fatal(err)
	}

	// get redis database env value
	redisDatabase := os.Getenv("REDIS_DATABASE")
	if redisDatabase == "" {
		err := errors.New("invalid REDIS_DATABASE env value")
		log.Fatal(err)
	}
	rd, err := strconv.Atoi(redisDatabase)
	if err != nil {
		log.Fatal(err)
	}

	// initialice message backend with Approve and Deliver methods
	backend, err := backendIot.NewBackend(providers)
	if err != nil {
		log.Fatal(err)
	}

	// initialize api
	svc, err := service.NewMessageService("IoT", service.ServiceConfig{
		Port: port,

		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisDatabase: rd,

		Approve: backend.Approve,
		Deliver: backend.Deliver,
	})
	if err != nil {
		log.Fatal(err)
	}

	// run message service
	err = svc.Run()
	if err != nil {
		log.Fatal(err)
	}
}
