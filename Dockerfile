FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -o app ./cmd/garden-equipment-log

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app ./garden-equipment-log

CMD ["./garden-equipment-log"]