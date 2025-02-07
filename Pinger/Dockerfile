FROM golang:1.22.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app .
CMD ["chmod 666 /var/run/docker.sock"]
CMD ["./main"]