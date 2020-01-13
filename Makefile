#
# SO variables
#
# DOCKER_USER
# DOCKER_PASS
#

#
# Internal variables
#
VERSION=0.2.0
SVC=messages-iot-api
BIN_PATH=$(PWD)/bin
BIN=$(BIN_PATH)/$(SVC)
REGISTRY_URL=$(DOCKER_USER)

#
# SVC variables
#
PORT=5060
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE=1
PROVIDERS=nats,mqtt
PROVIDER_NATS_API_KEY=123
PROVIDER_MQTT_API_KEY=456

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

clean-proto cp:
	@echo "[clean-proto] Cleaning proto files..."
	@rm -rf proto/*.pb.go || true

proto pro: clean-proto
	@echo "[proto] Generating proto file..."
	@protoc --go_out=plugins=grpc:. ./proto/*.proto 

run r:
	@echo "[running] Running service..."
	@PROVIDERS=$(PROVIDERS) PROVIDER_NATS_API_KEY=$(PROVIDER_NATS_API_KEY) PROVIDER_MQTT_API_KEY=$(PROVIDER_MQTT_API_KEY) PORT=$(PORT) REDIS_HOST=$(REDIS_HOST) REDIS_PORT=$(REDIS_PORT) REDIS_DATABASE=$(REDIS_DATABASE) ./bin/$(SVC)

build b:
	@echo "[build] Building service..."
	@cd cmd && go build -o $(BIN)

linux l: 
	@echo "[build-linux] Building service..."
	@cd cmd && GOOS=linux GOARCH=amd64 go build -o $(BIN)

docker d:
	@echo "[docker] Building image..."
	@docker build -t $(SVC):$(VERSION) .
	
docker-login dl:
	@echo "[docker] Login to docker..."
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)

push p: linux docker docker-login
	@echo "[docker] pushing $(REGISTRY_URL)/$(SVC):$(VERSION)"
	@docker tag $(SVC):$(VERSION) $(REGISTRY_URL)/$(SVC):$(VERSION)
	@docker push $(REGISTRY_URL)/$(SVC):$(VERSION)

compose co:
	@echo "[docker-compose] Running docker-compose..."
	@docker-compose build
	@docker-compose up

stop s: 
	@echo "[docker-compose] Stopping docker-compose..."
	@docker-compose down

.PHONY: clean c clean-proto cp proto pro run r build b linux l docker d docker-login dl push p compose co stop s