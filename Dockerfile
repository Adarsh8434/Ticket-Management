FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ticket-server ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/ticket-server .

EXPOSE 8080

CMD ["./ticket-server"]