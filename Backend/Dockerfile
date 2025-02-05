FROM golang:1.22.5

WORKDIR /backend

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o backend cmd/main.go

EXPOSE 8081

CMD ["./backend"]