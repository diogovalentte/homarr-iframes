FROM golang:1.21.6

WORKDIR /app

COPY . .

ENV GIN_MODE=release

WORKDIR /app

RUN go mod download

RUN go build -o main .

CMD ["./main"]
