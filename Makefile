GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

mock: 
	mockery --dir todo/repository --all --output todo/mocks/repository
	mockery --dir todo/service --all --output todo/mocks/service
run:
	air
test:
	go test ./...
mock-test:
	make mock
	make test
build:
	go build -o go-rengan cmds/app/main.go
.PHONY: test/cover
test/cover:
	mkdir -p coverage
	$(GOTEST) -v -coverprofile=coverage/coverage.out ./...
	$(GOCOVER) -func=coverage/coverage.out
	$(GOCOVER) -html=coverage/coverage.out -o coverage/coverage.html
wire:
	wire gen go-rengan/dep