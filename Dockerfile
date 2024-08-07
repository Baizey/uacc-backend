FROM golang:1.22 as builder
LABEL authors="Baizey"
WORKDIR /app

COPY . .
RUN go build -o server .

FROM ubuntu:latest as deploy
WORKDIR /app

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/* \

ENV PORT=8080

EXPOSE $PORT

COPY --from=builder /app/server .
CMD ./server