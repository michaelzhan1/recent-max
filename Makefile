.PHONY: generator
generator:
	docker compose up --build generator

.PHONY: server
server:
	docker compose up --build server