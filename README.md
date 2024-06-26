# Homarr iFrames

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/8df579cb-9cc9-4bad-a1da-f0cf015e741b)


An API that connects with multiple sources and creates an iFrame to be used in any dashboard, like [Homarr](https://github.com/ajnart/homarr).

The iFrames will be available under the API routes, like `/v1/iframe/linkwarden`. These routes also accept query parameters to change the iFrame, like limiting the number of items or specifying if you want the iFrames to check for updates automatically, and, if the source contents changes (like adding new bookmarks on Linkwarden), the iframe should reload.

- You can check all query parameters in the API docs.

# Sources

The API can create iFrames for multiple sources, like the **Vikunja** source that creates an iFrame with your tasks, or the **Linkwarden** source that creates an iFrame with your bookmarks.

Some sources require specific information to work, like a service address or credentials. You need to provide this information using environment variables.

The way you provide these environment variables depends on how you run the API.

- A list of the sources can be found [here](/docs/SOURCES.md).

# API docs

After starting the API, you can find the API docs under the path `/v1/swagger/index.html`, like `http://192.168.1.44/v1/swagger/index.html` or `https://sub.domain.com/v1/swagger/index.html`, depending on how you access the API.

# Notes

When you add an iFrame widget in your dashboard, it's **>your<** web browser that fetches the iFrame from the API and shows it to you, not your dashboard application. So your browser needs to be able to access the API, that's how an iFrame works.

- **Examples**:
  - If you run the API on your server, you need to use your server IP address + port, and you need to make sure your browser can access this IP + port.
  - If you're accessing your dashboard with a domain and using HTTPS, you also need to access this API with a domain and using HTTPS. If you try to use HTTP with your HTTPS, your browser will block the iFrame.

# How to run:

- **For Docker and Docker Compose**: by default, the API will be available on port `8080` and is not accessible by other machines. To be accessible by other machines, you need to run the API behind a reverse proxy or run the container in [host network mode](https://docs.docker.com/network/drivers/host/).

- You can change the API port using the environment variable `PORT`.
  - Depending on the port you choose, you need to run the container with user `root` instead of the user `1000` used in the examples and the `docker-compose.yml` file.

## Using Docker:

1. Run the latest version:

```sh
docker run --name homarr-iframes -p 8080:8080 -e VARIABLE_NAME=VARIABLE_VALUE -e VARIABLE_NAME=VARIABLE_VALUE ghcr.io/diogovalentte/homarr-iframes:latest
```

## Using Docker Compose:

1. There is a `docker-compose.yml` file in this repository. You can clone this repository to use this file or create one yourself.
2. Create a `.env` file with the environment variables you want to provide to the API. It should be like the `.env.example` file and be in the same directory as the `docker-compose.yml` file.
3. Start the container by running:

```sh
docker compose up
```

## Manually:

1. Install the dependencies:

```sh
go mod download
```
2. Export the environment variables.
3. Run:

```sh
go run main.go
```

# Adding to Homarr (or any other dashboard)

1. In your Homarr dashboard, click on **Enter edit mode -> Add a tile -> Widgets -> iFrame**.
2. Click to edit the new iFrame widget.
3. Add the API URL (`http://192.168.1.15:8080`) + the source path (`/v1/iframe/linkwarden`) + query parameters, like `http://192.168.1.15:8080/v1/iframe/linkwarden?collectionId=1&limit=3&theme=dark`.

# IMPORTANT!

- This project doesn't have any authentication system, so anyone who can access the API will be able to get all information from the API routes, like your Vikunja tasks, Linkwarden bookmarks, etc. You can add an authentication portal like [Authelia](https://github.com/authelia/authelia) or [Authentik](https://github.com/goauthentik/authentik) in front of the project to secure it, that's how I do it.
