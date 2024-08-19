FROM golang:1.23-alpine AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o rediscli

RUN ./rediscli completion bash > /tmp/completion.sh


FROM debian:bullseye-slim
WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/rediscli /usr/local/bin/rediscli
COPY --from=builder /tmp/completion.sh /etc/bash_completion.d/rediscli


ENTRYPOINT ["rediscli"]


CMD ["--help"]
