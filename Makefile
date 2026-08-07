default: all

SERVER:=./bin/server
SERVER_SRC:=$(shell find cmd internal -type f -name '*.go') go.mod go.sum

$(SERVER): $(SERVER_SRC)
	@mkdir -p bin
	go build -o $(SERVER) ./cmd

.PHONY: all
all: $(SERVER)
	./bin/server

.PHONY: clean
clean:
	rm -rf bin

.PHONY: frontend
frontend:
	cd web && npm run dev
