FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
MAILER_BINARY=mailApp
LOGGER_BINARY=loggerServiceApp
AUTH_BINARY=authApp
LISTENER_BINARY=listenerApp
FRONT_BINARY=frontEndApp
VERSION=1.3.3
GITHASH ?= $(shell git describe --match 'v[0-9]*' --tags --always)
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query "Account" --output text)
AWS_REGION ?= $(shell aws ec2 describe-availability-zones --output text --query 'AvailabilityZones[0].[RegionName]')

# up_swarm: initialize and deploy docker swarm
up_swarm:
	@echo "Init swarm"
	docker swarm init
	@echo "Deploy swarm"
	docker stack deploy -c swarm.yml myapp


## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail build_listener build_front_linux
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## build: builds the binaries
build: build_broker build_auth build_logger build_mail build_listener build_front_linux

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_front_linux: builds the broker binary as a linux executable
build_front_linux:
	@echo "Building front end linux binary..."
	cd ../front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_BINARY} ./cmd/web
	@echo "Done!"

## build_front_linux: builds the broker binary as a linux executable
build_front_osx:
	@echo "Building front end linux binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"


## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash=$(GITHASH) -X main.version=${VERSION}" -o ${BROKER_BINARY} ./cmd/api
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

## start: starts the front end
start: build_front_osx
	@echo "Starting front end"
	cd ../front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## DOCKER HUB

## BUILD DOCKER IMAGES

# build and tag images for all services - might need separted version tag
build_and_tag_images: build_and_tag_mailer build_and_tag_logger build_and_tag_listener build_and_tag_broker build_and_tag_authentication build_and_tag_front-end
	@echo "Building images for docker hub repo"

# build and tag the broker service
build_and_tag_broker:
	@echo "Build and tag broker"
	cd ../broker-service && docker build -f broker-service.dockerfile -t tfosorcim/broker-service:${VERSION} .
	docker image tag tfosorcim/broker-service:${VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/broker-service:${VERSION}

# build and tag the authentication service
build_and_tag_authentication:
	@echo "Build and tag authentication"
	cd ../authentication-service && docker build -f authentication-service.dockerfile -t tfosorcim/authentication-service:${VERSION} .
	docker image tag tfosorcim/authentication-service:${VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/authentication-service:${VERSION}

# build and tag the logger service
build_and_tag_logger:
	@echo "Build and tag logger"
	cd ../logger-service && docker build -f logger-service.dockerfile -t tfosorcim/logger-service:${VERSION} .
	docker image tag tfosorcim/logger-service:${VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/logger-service:${VERSION}


# build and tag the mailer service
build_and_tag_mailer:
	@echo "Build and tag mailer"
	cd ../mail-service && docker build -f mail-service.dockerfile -t tfosorcim/mailer-service:${VERSION} .
	docker image tag tfosorcim/mailer-service:${VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/mailer-service:${VERSION}


# build and tag the listener service
build_and_tag_listener:
	@echo "Build and tag listener"
	cd ../listener-service && docker build -f listener-service.dockerfile -t tfosorcim/listener-service:${VERSION} .
	docker image tag tfosorcim/listener-service:${VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/listener-service:${VERSION}


# build and tag the front-end service
build_and_tag_front-end:
	@echo "Build and tag front-end"
	cd ../front-end && docker build -f front-end.dockerfile -t tfosorcim/front-end:${VERSION} .
	docker image tag tfosorcim/front-end:${VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/front-end:${VERSION}


## PUSH DOCKER IMAGES TO DOCKER HUB

# push all service images to docker hub
push_images: push_broker push_authentication push_logger push_mailer push_listener push_frontend
	@echo "Pushing images to docker hub"

# push broker image
push_broker:
	@echo "Push broker to docker hub"
	cd ../broker-service && docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/broker-service:${VERSION}
# push authentication service
push_authentication:
	@echo "Push authentication to docker hub"
	cd ../authentication-service && docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/authentication-service:${VERSION}

# push logger service
push_logger:
	@echo "Push logger to docker hub"
	cd ../logger-service && docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/logger-service:${VERSION}

# push mailer service
push_mailer:
	@echo "Push mailer to docker hub"
	cd ../mail-service && docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/mailer-service:${VERSION}

# push listener service
push_listener:
	@echo "Push listener to docker hub"
	cd ../listener-service && docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/listener-service:${VERSION}

# push fron-end service
push_frontend:
	@echo "Push front-end to docker hub"
	cd ../front-end && docker image push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/experiment/front-end:${VERSION}

