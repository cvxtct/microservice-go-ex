FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
MAILER_BINARY=mailApp
LOGGER_BINARY=loggerServiceApp
AUTH_BINARY=authApp
LISTENER_BINARY=listenerApp
VERSION=1.0.0

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mailer binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAILER_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth  binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## SWARM RELATED COMMANDS

## BUILD DOCKER IMAGES

build_and_tag_images: build_and_tag_mailer build_and_tag_logger build_and_tag_listener build_and_tag_broker build_and_tag_authentication
	@echo "Building images for docker hub repo"

build_and_tag_broker:
	@echo "Build and tag broker"
	cd ../broker-service && docker build -f broker-service.dockerfile -t tfosorcim/broker-service:${VERSION} .

build_and_tag_authentication:
	@echo "Build and tag authentication"
	cd ../authentication-service && docker build -f authentication-service.dockerfile -t tfosorcim/authentication-service:${VERSION} .

build_and_tag_logger:
	@echo "Build and tag logger"
	cd ../logger-service && docker build -f logger-service.dockerfile -t tfosorcim/logger-service:${VERSION} .

build_and_tag_mailer:
	@echo "Build and tag mailer"
	cd ../mail-service && docker build -f mail-service.dockerfile -t tfosorcim/mailer-service:${VERSION} .

build_and_tag_listener:
	@echo "Build and tag listener"
	cd ../listener-service && docker build -f listener-service.dockerfile -t tfosorcim/listener-service:${VERSION} .

## PUSH DOCKER IMAGES TO DOCKER HUB

push_images: push_broker push_authentication push_logger push_mailer push_listener
	@echo "Pushing images to docker hub"

push_broker:
	@echo "Push broker to docker hub"
	cd ../broker-service && docker push tfosorcim/broker-service:${VERSION}

push_authentication:
	@echo "Push authentication to docker hub"
	cd ../authentication-service && docker push tfosorcim/authentication-service:${VERSION}

push_logger:
	@echo "Push logger to docker hub"
	cd ../logger-service && docker push tfosorcim/logger-service:${VERSION}

push_mailer:
	@echo "Push mailer to docker hub"
	cd ../mail-service && docker push tfosorcim/mailer-service:${VERSION}

push_listener:
	@echo "Push listener to docker hub"
	cd ../listener-service && docker push tfosorcim/listener-service:${VERSION}

