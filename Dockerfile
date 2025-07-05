# syntax=docker/dockerfile:1

FROM golang:1.24 AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/main.go

# Start fresh with a small image
FROM debian:bookworm-slim

WORKDIR /app

ENV DBURI=""
ENV ACCOUNTSID=""
ENV AUTHTOKEN=""
ENV FROMNUMBER=""
ENV CONTENTSID=""

COPY --from=build /app/main .

EXPOSE 5000

CMD ["./main"]
