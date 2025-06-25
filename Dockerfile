FROM golang:1.23.4 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o main .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

ENV GIN_MODE=release
ENV PORT=8080
ENV TZ=UTC

COPY --from=build /app/main .

CMD ["./main"]
