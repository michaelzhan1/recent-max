.PHONY: generator
generator:
	go build -o bin/generator cmd/generator/main.go
	./bin/generator

.PHONY: server
server:
	go build -o bin/server cmd/server/main.go
	./bin/server