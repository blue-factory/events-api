#
# SO variables
#
# DOCKER_USER
# DOCKER_PASS
#

#
# Internal variables
#
VERSION=0.1.2
SVC=messages-iot-api
BIN_PATH=$(PWD)/bin
BIN=$(BIN_PATH)/$(SVC)
REGISTRY_URL=$(DOCKER_USER)

#
# SVC variables
#
PORT=5060
HOST=localhost
MESSAGES_HOST=localhost
MESSAGES_PORT=5050
PROVIDERS=nats,mqtt
PROVIDER_NATS_API_KEY=<PROVIDER_NATS_API_KEY>
PROVIDER_MQTT_API_KEY=<PROVIDER_MQTT_API_KEY>

clean c:
	@echo "[clean] Cleaning bin folder..."
	@rm -rf bin/

run r:
	@echo "[running] Running service..."
	@PROVIDERS=$(PROVIDERS) PROVIDER_SENDGRID_API_KEY=$(PROVIDER_SENDGRID_API_KEY) MESSAGES_HOST=$(MESSAGES_HOST) MESSAGES_PORT=$(MESSAGES_PORT) HOST=$(HOST) PORT=$(PORT) PROVIDER_MANDRILL_API_KEY=$(PROVIDER_MANDRILL_API_KEY) PROVIDER_SES_AWS_KEY_ID=$(PROVIDER_SES_AWS_KEY_ID) PROVIDER_SES_AWS_SECRET_KEY=$(PROVIDER_SES_AWS_SECRET_KEY) PROVIDER_SES_AWS_REGION=$(PROVIDER_SES_AWS_REGION) go run cmd/main.go

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

.PHONY: clean c run r build b linux l docker d docker-login dl push p compose co stop s