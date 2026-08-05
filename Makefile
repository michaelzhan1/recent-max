default: server

.PHONY: server
server:
	docker compose up server

.PHONY: build
build:
	docker compose up --build

.PHONY: frontend
frontend:
	cd web && npm run dev
