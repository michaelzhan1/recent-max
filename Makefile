default: all

.PHONY: generator
generator:
	docker compose up --build generator

.PHONY: server
server:
	docker compose up --build server\

.PHONY: all
all:
	docker compose up --build