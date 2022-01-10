TAG := $(shell git describe --tags)

ROOT=$(shell pwd)
CMD= $(ROOT)/cmd
WEB= $(CMD)/web
WORKER= $(CMD)/worker
MIGRATE= $(CMD)/migrate

web-build:
	ROOT=$(ROOT) $(MAKE) -C $(WEB)

web-dev:
	ROOT=$(ROOT) air -c cmd/web/air.toml

worker-build:
	ROOT=$(ROOT) $(MAKE) -C $(WORKER)

worker-dev:
	ROOT=$(ROOT) $(MAKE) -C $(WORKER) dev

migrate-build:
	ROOT=$(ROOT) $(MAKE) -C $(MIGRATE) build

migrate-create:
	migrate create -ext sql -dir db/migrations $(name)

unit-test:
	go test ./... -v -cover

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml --port=7000