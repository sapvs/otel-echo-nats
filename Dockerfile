FROM golang:1.20-alpine3.18 AS builder
WORKDIR /app

COPY go.mod .
COPY go.sum . 
RUN go mod download

COPY . .

ARG DIR 
RUN go build -o main ${DIR}/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .

ENTRYPOINT ./main