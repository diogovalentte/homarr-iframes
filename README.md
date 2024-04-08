# Homarr iFrames

An API that gets data from multiple sources and creates a nice HTML code to be used in an iFrame (designed to be used in [Homarr](https://github.com/ajnart/homarr)).

The iFrames will be available under the API routes, like `/v1/iframes/linkwarden`. These routes also accept query parameters to change the iFrame HTML, like limiting the number of items or selecting the iFrame theme (light or dark) to match your Homarr dashboard theme.

# Sources

The API can create iFrames for multiple sources. Examples:

- The **Vikunja** source creates an iFrame with your tasks.
- The **Linkwarden** source creates an iFrame with your bookmarks.

Some sources require specific information to work, like a service address or credentials. You need to provide this information using environment variables.

The way you provide these environment variables depends on how you run the API.

- A list of the sources can be found [here](/docs/SOURCES.md).

# API docs

After starting the API, you can find the API docs under the path `/v1/swagger/index.html`, like `http://192.168.1.44/v1/swagger/index.html` or `https://sub.domain.com/v1/swagger/index.html`, depending on how you access the API.

# How to run:

## Using Docker:

1. Run the latest version:

```sh
docker run --name homarr-iframes -p 8080:8080 -e VARIABLE_NAME=VARIABLE_VALUE -e VARIABLE_NAME=VARIABLE_VALUE ghcr.io/diogovalentte/homarr-iframes:latest
```

## Using Docker Compose:

1. There is a `docker-compose.yml` file in this repository. You can clone this repository to use this file or create one yourself.
2. Create an `.env` file with the environment variables you want to provide to the API. It should be like the `.env.example` file and be in the same directory as the `docker-compose.yml` file.
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

# Adding to Homarr

1. In your Homarr dashboard, click on **Enter edit mode -> Add a tile -> Widgets -> iFrame**.
2. Click to edit the new iFrame widget.
3. Add the API URL (`http://192.168.1.15:8080`) + the source path (`/v1/iframe/linkwarden`) + query parameters, like `http://192.168.1.15:8080/v1/iframe/linkwarden?collectionId=1&limit=3&theme=dark`.

# IMPORTANT!

- This API doesn't have any authentication system, so anyone who can access the API will be able to get all information from all sources, like your Vikunja tasks, Linkwarden bookmarks, etc. You can add an authentication portal like [Authelia](https://github.com/authelia/authelia) or [Authentik](https://github.com/goauthentik/authentik) in front of the API to secure it, this is how I do it.
