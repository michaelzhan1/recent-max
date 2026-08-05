FROM golang:1.26 AS builder

WORKDIR /app

COPY . .
RUN go mod download

FROM builder AS server-builder

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /out/server ./cmd/server

FROM debian:bookworm-slim AS server

COPY --from=server-builder /out/server /usr/local/bin/server

ENTRYPOINT ["/usr/local/bin/server"]
