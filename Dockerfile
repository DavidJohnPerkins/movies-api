# syntax=docker/dockerfile:1

FROM golang:1.26.1 AS builder
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Final minimal image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 8080
CMD ["./app"]
