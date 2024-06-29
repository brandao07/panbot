# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
BINARY_NAME=panbot
BINARY_PATH=./build/panbot
PANBOT_MAIN_PATH=./cmd/panbot/main
TMP_PATH=./tmp

# Docker parameters
DOCKER_COMPOSE=docker-compose

# Air parameters
AIR=air
AIR_CONFIG=.air.toml

all: build

build:
	$(GOBUILD) -o $(BINARY_PATH)/$(BINARY_NAME) $(PANBOT_MAIN_PATH)/main.go

run:
	$(GOCMD) run $(PANBOT_MAIN_PATH)/main.go

clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)/$(BINARY_NAME)
	rm -rf ${TMP_PATH}

mod-download:
	$(GOMOD) download

docker-up:
	$(DOCKER_COMPOSE) up

docker-down:
	$(DOCKER_COMPOSE) down

docker-rebuild: docker-down docker-up

air:
	$(AIR) -c $(AIR_CONFIG)

.PHONY: all build run clean mod-download docker-build docker-up docker-down docker-rebuild air
