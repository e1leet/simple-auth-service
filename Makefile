CONFIG_PATH = ./configs/config.local.yaml
POSTGRES_URI = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

.PHONY: run
run:
	go run ./cmd/app/main.go --config=$(CONFIG_PATH)

.PHONY: lint
lint:
	golangci-lint run ./... --config=./.golangci.yaml

.PHONY: test
test:
	go test -v --race --timeout=5m --cover ./...

.PHONY: swagger
swagger:
	swag init -g ./cmd/app/main.go

.PHONY: migrate
migrate:
	migrate -path migrations/ -database $(POSTGRES_URI) -version $(version)

.PHONY: migrate.up
migrate.up:
	migrate -path migrations/ -database $(POSTGRES_URI) up

.PHONY: migrate.down
migrate.down:
	migrate -path migrations/ -database $(POSTGRES_URI) down

