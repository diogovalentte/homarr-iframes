# Build stage
FROM golang:1.23.4 AS build

WORKDIR /app

# CGO (C bindings for Go)
# When CGO_ENABLED=0, Go builds a fully statically linked binary, meaning it does not rely on any external C libraries.
# This ensures that the binary is self-contained and portable.
ARG CGO_ENABLED=0

COPY . .

RUN go mod download

RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /app

ENV GIN_MODE=release
ENV PORT=8080
ENV TZ=UTC

RUN apk update && apk add --no-cache curl

HEALTHCHECK --interval=10s --timeout=10s --start-period=15s --retries=3 \
  CMD sh -c 'curl -f http://localhost:$PORT/v1/health | grep OK || exit 1'

# Copy the compiled binary from the build stage
COPY --from=build /app/main .

CMD ["./main"]
