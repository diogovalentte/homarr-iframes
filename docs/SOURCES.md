# Sources

- Each source is a route of the API that returns an iFrame.
- Some sources need some environment variables to work.
- Most sources have two environment variables for the address, one for the address that will be used in the iFrame links and another for the address that will be used in the API to get the data (`INTERNAL_`). If you don't provide the second one, the first one will be used in the API too.
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
- `INTERNAL_LINKWARDEN_ADDRESS`: your Linkwarden instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `LINKWARDEN_TOKEN`: an access token used to access your Linkwarden instance API to get your links. You can get it in **Settings -> Access Tokens -> New Access Token**.
- `LINKWARDEN_BACKGROUND_IMG_URL`: an image URL to be used as the background of each bookmark card.

# Vikunja

This source creates an iFrame with links to the tasks from your [Vikunja](https://github.com/go-vikunja/vikunja) instance.

- It automatically sorts the tasks by **due date** (ascendent), **end date** (ascendent), and **created date** (descendent), and also filters to return only tasks that are **not done**.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/787ff13a-a81f-42b4-a3a4-9f0892ca815f)

To use this source, you'll need to provide the following environment variables:

- `VIKUNJA_ADDRESS`: your Vikunja instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_VIKUNJA_ADDRESS`: your Vikunja instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `VIKUNJA_TOKEN`: an access token used to access your Vikunja instance API to get your tasks. You can get it by going to **Settings -> API Tokens -> Create a Token -> In "Tasks", select "Read One" and "Read All"; In "Projects", select "Read One" and "Read All" -> Create Token**.
  - If you want to add a button to set the task as done in the iframe, add the permission **Update**.
- `VIKUNJA_BACKGROUND_IMG_URL`: an image URL to be used as the background of each task card.

# Media Requests

This source creates an iFrame with your media requests from your [Overseerr](https://github.com/sct/overseerr) and [Jellyseerr](https://github.com/Fallenbagel/jellyseerr) instances.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7f374beb-e392-4ee9-94fc-4d1556f65e7c)

To use this source, you'll need to provide the following environment variables:

- `OVERSEERR_ADDRESS`: your Overseerr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`. It'll be used in the links in the iframe. If `INTERNAL_OVERSEERR_ADDRESS` is not provided, it'll also be used by the API to get the data from Overseerr.
- `INTERNAL_OVERSEERR_ADDRESS`: your Overseerr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`. It'll be used by the API to get the data from Overseerr.
- `OVERSEERR_API_KEY`: an API key used to access your Overseerr instance API to get your media requests. You can get it by going to **Settings -> General -> API Key**.

- `JELLYSEERR_ADDRESS`: your Jellyseerr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`. It'll be used in the links in the iframe. If `INTERNAL_JELLYSEERR_ADDRESS` is not provided, it'll also be used by the API to get the data from Jellyseerr.
- `INTERNAL_JELLYSEERR_ADDRESS`: your Jellyseerr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`. It'll be used by the API to get the data from Jellyseerr.
- `JELLYSEERR_API_KEY`: an API key used to access your Jellyseerr instance API to get your media requests. You can get it by going to **Settings -> General -> API Key**.

# Media Releases

This source creates an iFrame with media that is released today. There is also an indicator of whether the media is downloaded or not (available or not available).

- It gets the media from [Sonarr](https://github.com/Sonarr/Sonarr) and [Radarr](https://github.com/Radarr/Radarr).
- Set the container timezone to get a more precise result.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/17665aaf-5eb5-44df-a0db-9b99a5fe77fb)

To use this source, you'll need to provide the environment variables below, but you don't need to provide all of them, you can specify only the Sonarr variables for example.

- `SONARR_ADDRESS`: your Sonarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_SONARR_ADDRESS`: your Sonarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `SONARR_API_KEY`: an access API key used to access your Sonarr instance API to get your media. You can get it by going to **Settings -> General -> API Key**.

- `RADARR_ADDRESS`: your Radarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_RADARR_ADDRESS`: your Radarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `RADARR_API_KEY`: an access API key used to access your Radarr instance API to get your media. You can get it by going to **Settings -> General -> API Key**.

# Uptime Kuma

This source creates an iFrame with the number of UP and DOWN sites from a [Uptime Kuma]() status page.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7b0e2cfc-2edc-41d4-9551-72df189591d4)

To use this source, you'll need to provide the following environment variables:

- `UPTIMEKUMA_ADDRESS`: your Uptime Kuma instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.

# Cinemark Brasil

This source gets on display movies of specific Cinemark theaters (only in Brazil) and creates an iFrame. It shows some info about the films and has links to their pages.

- You have to specify which theaters to get movies from. I recommend specifying all theaters in your city.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7071b022-fe90-4db7-874b-8b88d0298641)

# Alarms

This source shows **alarms** (warnings, errors, failed jobs notifications, etc.) from multiple services in one central place.

![image](https://github.com/user-attachments/assets/15e26b24-8d4b-4243-b239-e6f4c5056712)

To use this source, you must provide environment variables for each service from which you want to show alarms. You also need to specify the services' names in the iframe URL query parameter `alarms`.

Below are the available services that you can use in this iframe and the required environment variables:

## Netdata

Shows [Netdata](https://github.com/netdata/netdata) alerts, like high RAM/CPU usage alerts.

- `NETDATA_ADDRESS`: your Netdata instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_NETDATA_ADDRESS`: your Netdata instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `NETDATA_TOKEN`: an access token used to access your Netdata instance API to get your alarms. See how to get it [here](https://learn.netdata.cloud/docs/netdata-cloud/authentication-&-authorization/api-tokens).

## Radarr, Sonarr, and Prowlarr

Shows health messages from your [Sonarr](https://github.com/Sonarr/Sonarr), [Radarr](https://github.com/Radarr/Radarr), and [Prowlarr](https://github.com/Prowlarr/Prowlarr) instances, like when an index fails or Sonarr can't connect to a download client.

- `SONARR_ADDRESS`: your Sonarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_SONARR_ADDRESS`: your Sonarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `SONARR_API_KEY`: an access API key used to access your Sonarr instance API to get your media. You can get it by going to **Settings -> General -> API Key**.

- `RADARR_ADDRESS`: your Radarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_RADARR_ADDRESS`: your Radarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `RADARR_API_KEY`: an access API key used to access your Radarr instance API to get your media. You can get it by going to **Settings -> General -> API Key**.

- `PROWLARR_ADDRESS`: your Prowlarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_PROWLARR_ADDRESS`: your Prowlarr instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `PROWLARR_API_KEY`: an access API key used to access your Prowlarr instance API. You can get it by going to **Settings -> General -> API Key**.

## Speedtest Tracker

Shows a warning if the last speed test from your [Speedtest Tracker](https://github.com/alexjustesen/speedtest-tracker) instance failed.

- `SPEEDTEST_TRACKER_ADDRESS`: your Speedtest Tracker instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_SPEEDTEST_TRACKER_ADDRESS`: your Speedtest Tracker instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `SPEEDTEST_TRACKER_TOKEN`: your API token used to access your Speedtest Tracker instance. You can get it by going to **Settings -> API Tokens -> Create API Token button**.

## Pi-hole

Shows [Pi-hole](https://github.com/pi-hole/pi-hole) diagnostic messages, like a high load.

- `PIHOLE_ADDRESS`: your Pihole instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_PIHOLE_ADDRESS`: your Pihole instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `PIHOLE_TOKEN`: an access API key used to access your Pihole instance API. You can get it by going to **Settings -> API -> Show API Token button**.

## Kavita

Shows your [Kavita](https://github.com/Kareadita/Kavita) instance [media issues](https://wiki.kavitareader.com/troubleshooting/media-errors) that Kavita detects when analyzing your media, like corrupted files.

- `KAVITA_ADDRESS`: your Kavita instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_KAVITA_ADDRESS`: your Kavita instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `KAVITA_USERNAME`: your Kavita username.
- `KAVITA_PASSWORD`: your Kavita password.

## Kaizoku

Shows warnings if there are failed jobs in your [Kaizoku](https://github.com/oae/kaizoku) queues.

- `KAIZOKU_ADDRESS`: your Kaizoku instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_KAIZOKU_ADDRESS`: your Kaizoku instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.

## ChangeDetection.io

Shows cards for your watches' errors and changes from your [ChangeDetection.io](https://github.com/dgtlmoon/changedetection.io) instance.

- `CHANGEDETECTIONIO_ADDRESS`: your ChangeDetection.io instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_CHANGEDETECTIONIO_ADDRESS`: your ChangeDetection.io instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `CHANGEDETECTIONIO_API_KEY`: an access API key used to access your ChangeDetection.io instance API. You can get it by going to **Settings -> API -> Generate API Key button**.
- `CHANGEDETECTIONIO_CHANGED_LAST_HOURS`: number of hours to indicate if the iframe should show a watch change. If the watch's last changed time is within the last `x` hours, it'll show the watch, else no. Defaults to 24.
