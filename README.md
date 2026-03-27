# Homarr iFrames

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/8df579cb-9cc9-4bad-a1da-f0cf015e741b)

This project connects to multiple self-hosted applications (referred to as sources) and exposes their data through embeddable iFrames.

Despite the name, the generated iFrames can be used in any dashboard, not only [Homarr](https://github.com/ajnart/homarr).

Each source is exposed through an API route such as:

```
/v1/iframe/linkwarden
```

These routes accept query parameters that allow you to:

- Limit the number of displayed items.
- Enable automatic update checks (the iFrame reloads if the source data changes).
- Customize appearance (when supported by the source).

All available query parameters are documented in the API documentation.

---

# Sources

The API supports multiple sources. Examples:

- Vikunja — displays your tasks.
- Linkwarden — displays your bookmarks.

Each source may require specific environment variables, such as:

- Application URL
- API tokens
- Credentials

How you provide these variables depends on how you run the API (Docker, Docker Compose, or manually).

A complete list of supported sources is available [here](/docs/SOURCES.md):

# API Documentation

Swagger documentation is available at:

```
/v1/swagger/index.html
```

Examples:

- `http://192.168.1.44/v1/swagger/index.html`
- `https://sub.domain.com/v1/swagger/index.html`

The exact URL depends on how and where the API is hosted.

---

# Notes

## Adding the iFrame to Homarr

1. Enter edit mode.
2. Add a new item.
3. Select the iFrame item and add it.
4. Configure the widget with:

```
<API_URL>/v1/iframe/<source>?<query_parameters>
```

Example:

```
http://192.168.1.15:8080/v1/iframe/linkwarden?collectionId=1&limit=3&theme=dark
```

## How iFrame Access Works

When you add the iFrame to your dashboard:

- Your browser requests the iFrame directly from this API.
- The dashboard server does not proxy the request.

This means:

- Your browser must be able to access the API.
- The protocol must match.

Examples:

- If the API runs on port 5000, access it via `http://SERVER_IP:5000`.
- If your dashboard uses HTTPS, the API must also be served over HTTPS.
  Mixing HTTP and HTTPS will cause the browser to block the iFrame.

## No Built-in Authentication

This project does not provide authentication.

Anyone who can access the API can retrieve data from all configured sources.

To secure it, place an authentication layer in front of the API, such as:

- [Authelia](https://github.com/authelia/authelia)
- [Authentik](https://github.com/goauthentik/authentik)

---

# Running

## Docker and Docker Compose

By default:

- The API runs on port `8080`.
- It is not accessible externally unless configured.

To make it accessible from other machines:

- Run it behind a reverse proxy, or
- Use [host network mode](https://docs.docker.com/network/drivers/host/).

You can change the port using the `PORT` environment variable.

## Using Docker

1. Run the latest image:

```sh
docker run \
  --name homarr-iframes \
  -p 8080:8080 \
  -e VARIABLE_NAME=VARIABLE_VALUE \
  ghcr.io/diogovalentte/homarr-iframes:latest
```

## Using Docker Compose

1. Use the provided `docker-compose.yml` or create your own.
2. Create a `.env` file (based on `.env.example`) in the same directory.
3. Start the service:

```sh
docker compose up
```

## Running Manually

1. Install dependencies:

```sh
go mod download
```

2. Export the required environment variables.
3. Run the API:

```sh
go run main.go
```
