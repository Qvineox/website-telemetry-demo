FROM golang:1.22.2-alpine3.19 AS builder
LABEL authors="lysak.yaroslav00@yandex.ru"

# ref: https://habr.com/ru/companies/otus/articles/660301/

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux

WORKDIR /src

COPY ./go.mod /src/go.mod
RUN go mod download

COPY . /src
RUN go mod tidy

RUN mkdir /app_build
RUN go build -ldflags="-s -w" -o /app_build/go_build.exe ./cmd

FROM alpine:3.19.1 as web_server

COPY ./web ./web
COPY --from=builder /app_build .

CMD ["./go_build.exe"]