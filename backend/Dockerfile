# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
