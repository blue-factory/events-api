# Messages IoT API

Microservice implemented from [Messages Lib](https://github.com/microapis/messages-lib) is responsible for sending iot messages through some real time providers such as NATS or MQTT.

As explained in the Messages Lib repository, it can be seen that there are three models, messages, channel and provider. To know more you can read the readme of Messages Lib.

## Channels

Corresponds to the type of notification that could be sent, for this there must be the implementation of that channel as a gRPC service through an API.

## Providers

The provider is an attribute of **channel** and allows to identify what types of messages are available for a specific channel.

In this api we will find the implementation of only 3 providers. Find **Providers** implementation at the [`./provider`](./lib) folder.

| Name                                   | Protocol         | Description          | ENV (each with prefix `PROVIDER_`) |
| -------------------------------------- | ---------------- | -------------------- | ---------------------------------- |
| [Nats](https://nats.io/)               | nats             | On premise its free. | `NATS_API_KEY`: string             |
| [MQTT](https://mqtt.org/)              | mqtt             | On premise its free. | `MQTT_API_KEY`: string             |
| [AWS IoT](https://aws.amazon.com/iot/) | mqtt             | Coming soon.         | Coming soon.                       |
| [Azure IoT](https://mandrill.com/)     | mqtt, amqp, http | Coming soon.         | Coming soon.                       |

## gRPC Service

```go
service MessageBackendService {
  rpc Approve(MessageBackendApproveRequest) returns (MessageBackendApproveResponse) {}
  rpc Deliver(MessageBackendDeliverRequest) returns (MessageBackendDeliverResponse) {}
}
```

## Model

```go
IoT {
  topic:      string
  payload:    map[string]string
  Provider:   string    // nats, mqtt
}
```

## Commands (Development)

`make build`: build user service for osx.

`make linux`: build user service for linux os.

`make docker`: build docker.

`docker run -it -p 5060:5060 messages-iot-api`: run docker.

**Run Iot Messages API:**

```sh
HOST=<host> \
PORT=<port> \
MESSAGES_HOST=<messages-host> \
MESSAGES_PORT=<messages-port> \
PROVIDERS=<providers> \
PROVIDER_NATS_API_KEY=<> \
PROVIDER_MQTT_API_KEY=<> \
./bin/messages-iot-api
```

## TODO

- [ ] Task 1.
- [ ] Task 2.
- [ ] Task 3.
