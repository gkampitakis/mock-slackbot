APP_NAME = mock-slackbot
BUILD_DIR = ${PWD}/build
EXECUTABLES = nodemon golangci-lint
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(warning "No $(exec) in PATH")))

run: build
	./build/$(APP_NAME)

dev:
	nodemon --exec go run *.go --signal SIGTERM

dependencies:
	go mod download
	go install github.com/mfridman/tparse@latest

clean:
	rm -rf ./build

build: clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

lint: 
	golangci-lint run -c ./golangci.yml

format:
	golines -w .

test: 
	go test ./... -cover -v -count=1 -json | tparse -all

test-cov:
	go test ./... -count=1 -coverprofile=coverage.out && \
	go tool cover -html=coverage.out
