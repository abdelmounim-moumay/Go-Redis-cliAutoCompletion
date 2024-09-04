<<<<<<< HEAD
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
=======
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


CMD ["bash", "--init-file", "/etc/bash_completion.d/rediscli"]
>>>>>>> e3476a6ccf315f8b7a4b2102ad12dfe1acb61df1
