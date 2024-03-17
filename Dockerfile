FROM golang:latest AS builder

COPY . /app

WORKDIR /app

RUN go build -o main .

FROM debian:latest

COPY --from=builder /app/main /usr/local/bin/main

ENV PORT 10011
EXPOSE 10011

CMD ["main"]

