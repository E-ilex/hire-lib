FROM golang:1.21-alpine AS build

RUN apk add --no-cache gcc g++ git openssh-client

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 3000

ENV CGO_ENABLED=1

CMD ["./main"]
