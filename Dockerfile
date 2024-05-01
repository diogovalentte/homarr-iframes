FROM golang:1.21.6

WORKDIR /app

COPY . .

ENV GIN_MODE=release
ENV PORT=8080
ENV TZ=UTC

WORKDIR /app

RUN go mod download

RUN go build -o main .

CMD ["./main"]
