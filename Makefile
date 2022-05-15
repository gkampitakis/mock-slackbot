APP_NAME = slack-bot
BUILD_DIR = ${PWD}/build
EXECUTABLES = nodemon golangci-lint
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(warning "No $(exec) in PATH")))

run: build
	./build/$(APP_NAME)

dev:
	nodemon --exec go run *.go --signal SIGTERM

depencies:
	go mod download

clean:
	rm -rf ./build

build: clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

lint: 
	golangci-lint run -c ./golangci.yml

format:
	golines -w .
