# Sources

- Each source has an API route to return an iFrame.
- Some sources need some environment variables to work, if you do not specify them, the source will not work, and when you try to request this source, it'll return an error.
- Some sources need some query arguments to work, you can check the [API docs](https://github.com/diogovalentte/homarr-iframes/tree/main?tab=readme-ov-file#api-docs) to see which arguments are obligatory.
- This API doesn't have any authentication system, so anyone who can access the API will be able to get all information from all sources, like your Vikunja tasks, Linkwarden bookmarks, etc. You can add an authentication portal like Authelia or [Authentik](https://github.com/goauthentik/authentik) in front of the API to secure it, this is how I do it.
- Some iFrames display date information, set the Docker container timezone to get a better result.
- The iFrames design is based on the Homarr widget to show media requests from apps like [Jellyseerr](https://github.com/Fallenbagel/jellyseerr) and [Overseerr](https://github.com/sct/overseerr):

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/9083c67a-9bbf-4430-8ba9-929cd9b0d0ab)

---

# Linkwarden

This source creates an iFrame with your bookmarks from your [Linkwarden](https://github.com/linkwarden/linkwarden) instance. It has links to the bookmark link and the bookmark Linkwarden collection.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/90271b2c-dc4f-4ee7-a6d3-f256e12cad81)

To use this source, you'll need to provide the following environment variables:

- `LINKWARDEN_ADDRESS`: your Linkwarden instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `LINKWARDEN_TOKEN`: an access token used to access your Linkwarden instance API to get your links. You can get it by going to **Settings -> Access Tokens -> New Access Token**.

# Vikunja

This source creates an iFrame with your tasks from your [Vikunja](https://github.com/go-vikunja/vikunja) instance. It has links to the tasks.

- It automatically sorts the tasks by **due date** (ascendent), **end date** (ascendent), and **created date** (descendent), and also filters to return only tasks that are **not done**.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/787ff13a-a81f-42b4-a3a4-9f0892ca815f)

To use this source, you'll need to provide the following environment variables:

- `VIKUNJA_ADDRESS`: your Vikunja instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `VIKUNJA_TOKEN`: an access token used to access your Vikunja instance API to get your tasks. You can get it by going to **Settings -> API Tokens -> Create a Token -> In "Tasks", select "Read One" and "Read All"; In "Projects", select "Read One" and "Read All" -> Create Token**.
  - If you want to add a button to set the task as done in the iframe, also add the permission **Update**.

# Uptime Kuma

This source creates an iFrame with the number of UP and DOWN sites from a [Uptime Kuma]() status page.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7b0e2cfc-2edc-41d4-9551-72df189591d4)

To use this source, you'll need to provide the following environment variables:

- `UPTIMEKUMA_ADDRESS`: your Uptime Kuma instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.

# Cinemark Brasil

This source gets on display movies of specific Cinemark theaters (only in Brazil) and creates an iFrame. It shows some info about the films and has links to their pages.
- You have to specify which theaters to get movies from. I recommend you specifing all theaters in your city.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7071b022-fe90-4db7-874b-8b88d0298641)
