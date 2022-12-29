CONFIG_PATH = ./configs/config.local.yaml

.PHONY: run
run:
	go run ./cmd/app/main.go --config=$(CONFIG_PATH)

.PHONY: lint
lint:
	golangci-lint run ./... --config=./.golangci.yaml

.PHONY: test
test:
	go test -v --race --timeout=5m --cover ./...