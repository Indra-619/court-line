# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/main .

RUN chown -R app:app /app

USER app

EXPOSE 8080

CMD ["./main"]
