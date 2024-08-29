# Étape de build
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . /app
RUN go mod tidy
RUN go build -o rediscli

RUN ./rediscli completion bash > /tmp/completion.sh


FROM debian:bullseye-slim
WORKDIR /app


COPY --from=builder /app/rediscli /app/rediscli
COPY --from=builder /tmp/completion.sh /etc/bash_completion.d/rediscli


COPY config/config.json /app/config/config.json


RUN chmod +x /app/rediscli

ENTRYPOINT ["./rediscli"]
