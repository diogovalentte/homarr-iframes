# Sources

- Each **source** corresponds to an API route that returns an iFrame.
- Some sources require environment variables to function.

Most sources define two address variables:

- A public address (used in clickable links inside the iFrame)
- An internal address (prefixed with `INTERNAL_`) used by this project to fetch data

If the internal address is not provided, the public address is used for both.

**Example**

You access a service via `service.com`, but it is protected by an authentication proxy.

You would configure:

- Public address → `service.com`
- Internal address → Docker container hostname or internal network address

This allows the API to fetch data directly without going through the authentication layer.

---

# Query Parameters

Many sources support URL query parameters that modify behavior and appearance. Some sources require them.

See the [API documentation](https://github.com/diogovalentte/homarr-iframes/tree/main?tab=readme-ov-file#api-documentation).

# Security Notice

This project has no built-in authentication.

Anyone who can access the API can read all exposed data (tasks, bookmarks, media info, etc.). You should place an authentication gateway (e.g., [Authelia](https://github.com/authelia/authelia) or [Authentik](https://github.com/goauthentik/authentik)) in front of the API.

# Timezone

Some iFrames display dates. Set the Docker container timezone to match your system for accurate results.

---

# Linkwarden

Displays bookmarks from a [Linkwarden](https://github.com/linkwarden/linkwarden) instance, including:

- A link to the original bookmark
- A link to the bookmark collection inside Linkwarden

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/90271b2c-dc4f-4ee7-a6d3-f256e12cad81)

**Environment variables**

- `LINKWARDEN_ADDRESS`: your Linkwarden instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `INTERNAL_LINKWARDEN_ADDRESS`: your Linkwarden instance address, like `https://sub.domain.com` or `http://192.168.1.45:8080`.
- `LINKWARDEN_TOKEN`: an access token used to access your Linkwarden instance API to get your links. You can get it in **Settings -> Access Tokens -> New Access Token**.
- `LINKWARDEN_BACKGROUND_IMG_URL`: an image URL to be used as the background of each bookmark card.

# Vikunja

Displays tasks from a [Vikunja](https://github.com/go-vikunja/vikunja) instance.

- It automatically sorts the tasks by **due date** (ascending), **end date** (ascending), and **created date** (descending), and also filters to return only tasks that are **not done**.

Tasks are automatically:

- Filtered to only **not done**
- Sorted by:
  - Due date (ascending)
  - End date (ascending)
  - Created date (descending)

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/787ff13a-a81f-42b4-a3a4-9f0892ca815f)

**Environment variables**

- `VIKUNJA_ADDRESS`
- `INTERNAL_VIKUNJA_ADDRESS`
- `VIKUNJA_TOKEN`

Token permissions:

- Tasks → Read One, Read All
- Projects → Read One, Read All

Optional:

- Add **Update** permission to allow a “mark as done” button in the iFrame
- `VIKUNJA_BACKGROUND_IMG_URL` — background image URL for task cards

# Media Requests

Displays media requests from:
- [Overseerr](https://github.com/sct/overseerr)
- [Jellyseerr](https://github.com/Fallenbagel/jellyseerr)

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7f374beb-e392-4ee9-94fc-4d1556f65e7c)

**Overseerr variables**

- `OVERSEERR_ADDRESS`
- `INTERNAL_OVERSEERR_ADDRESS`
- `OVERSEERR_API_KEY`: (Settings -> General -> API Key).

**Jellyseerr variables**

- `JELLYSEERR_ADDRESS`
- `INTERNAL_JELLYSEERR_ADDRESS`
- `JELLYSEERR_API_KEY`: (Settings → General → API Key)

# Media Releases

Shows media releasing today and whether it is downloaded. For Lidarr, it shows how many tracks of an album are downloaded.

Sources:

- [Sonarr](https://github.com/Sonarr/Sonarr)
- [Radarr](https://github.com/Radarr/Radarr)
- [Lidarr](https://github.com/Lidarr/Lidarr)

Use the same timezone as your media containers for best results.

You may configure only the services you use.

![image](https://github.com/user-attachments/assets/461249d2-7979-47bd-913e-2247c31c8e2e)

**Environment variables**

- `SONARR_ADDRESS`
- `INTERNAL_SONARR_ADDRESS`
- `SONARR_API_KEY`: (API keys: Settings → General → API Key)

- `RADARR_ADDRESS`
- `INTERNAL_RADARR_ADDRESS`
- `RADARR_API_KEY`: (API keys: Settings → General → API Key)

- `LIDARR_ADDRESS`
- `INTERNAL_LIDARR_ADDRESS`
- `LIDARR_API_KEY`: (API keys: Settings → General → API Key)

# Uptime Kuma

Displays the number of UP and DOWN monitors from an [Uptime Kuma](https://github.com/louislam/uptime-kuma) status page.

![image](https://github.com/diogovalentte/homarr-iframes/assets/49578155/7b0e2cfc-2edc-41d4-9551-72df189591d4)

**Environment variables**

- `UPTIMEKUMA_ADDRESS`

# Cinemark Brasil

Displays currently showing movies for selected Cinemark theaters in Brazil and links to their pages.

You must specify which theaters to fetch (recommended: all theaters in your city).

![image](https://github.com/user-attachments/assets/aafe4a96-8b48-471d-8046-a189492e4137)

# Alarms

Aggregates alerts and warnings from multiple services into one dashboard view.

![image](https://github.com/user-attachments/assets/15e26b24-8d4b-4243-b239-e6f4c5056712)

You must:

1. Configure environment variables for each service.
2. Pass service names in the iFrame query parameter:

```
alarms=<service1,service2,...>
```

## Regex Filtering

You can filter alarms using the `ALARMS_REGEX` environment variable.

The regex matches a concatenated string (*without spaces*):

```
source summary URL status property value
```

**Example**: "NetdataSystem requires reboot after package updateshttps://netdata.domain.comWARNINGOS / System1 status" for the alarm below:

![image](https://github.com/user-attachments/assets/fbfc8053-e688-40e0-82fd-7be6a224cf89)

**Query parameter**

```
regex_include=true|false
```

- `true` → show only matching alarms (default)
- `false` → hide matching alarms

## Netdata

Shows alerts (CPU, RAM, etc.) from [Netdata](https://github.com/netdata/netdata)

- `NETDATA_ADDRESS`
- `INTERNAL_NETDATA_ADDRESS`
- `NETDATA_TOKEN`: see how to get it [here](https://learn.netdata.cloud/docs/netdata-cloud/authentication-&-authorization/api-tokens).

## Sonarr / Radarr / Lidarr / Prowlarr Health

Shows health warnings such as indexer failures or download client connection problems.

Each requires:

- `<SERVICE>_ADDRESS`
- `INTERNAL_<SERVICE>_ADDRESS`
- `<SERVICE>_API_KEY`: (API keys: Settings → General → API Key)

(Sonarr, Radarr, Lidarr, Prowlarr)

## Speedtest Tracker

Warns when the last speed test failed from your [Speedtest Tracker](https://github.com/alexjustesen/speedtest-tracker) instance.

- `SPEEDTEST_TRACKER_ADDRESS`
- `INTERNAL_SPEEDTEST_TRACKER_ADDRESS`
- `SPEEDTEST_TRACKER_TOKEN`: (API token: Settings -> API Tokens -> Create API Token button)

## Pi-hole

Displays [Pi-hole](https://github.com/pi-hole/pi-hole) diagnostic messages.

- `PIHOLE_ADDRESS`
- `INTERNAL_PIHOLE_ADDRESS`
- `PIHOLE_PASSWORD`: a password to access your Pi-hole instance API in Pi-hole versions after `v6.0`. It can be the password you use to log in to the Pi-hole interface, but I recommend using the **app password**, as it's the only one that works if you enable **2FA**. You can get the app password on **Settings -> Web interface / API**. Make sure you're on the **Expert** mode, and click on **Configure app password**.
- `PIHOLE_TOKEN`: a token to access your Pi-hole instance API in Pi-hole versions previous to v6.0. You can get it by going to **Settings -> API -> Show API Token button**.

## Kavita

Shows [media issues](https://wiki.kavitareader.com/troubleshooting/media-errors) detected by [Kavita](https://github.com/Kareadita/Kavita) (e.g., corrupted files).

- `KAVITA_ADDRESS`
- `INTERNAL_KAVITA_ADDRESS`
- `KAVITA_USERNAME`
- `KAVITA_PASSWORD`

## Kaizoku

Shows failed job warnings from [Kaizoku](https://github.com/oae/kaizoku).

- `KAIZOKU_ADDRESS`
- `INTERNAL_KAIZOKU_ADDRESS`

## ChangeDetection.io

Shows errors and detected page changes.s from your [ChangeDetection.io](https://github.com/dgtlmoon/changedetection.io) instance.

- `CHANGEDETECTIONIO_ADDRESS`
- `INTERNAL_CHANGEDETECTIONIO_ADDRESS`
- `CHANGEDETECTIONIO_API_KEY`: (API key: Settings -> API -> Generate API Key button)
- `CHANGEDETECTIONIO_CHANGED_LAST_HOURS`: If the watch's last changed time is within the last `x` hours, it'll show the watch. Defaults to 24.

## Backrest

Displays backup plans with warning or error status in the last 24 hours from your [Backrest](https://github.com/garethgeorge/backrest) instance.

- `BACKREST_ADDRESS`
- `INTERNAL_BACKREST_ADDRESS`
- `BACKREST_USERNAME`
- `BACKREST_PASSWORD`

## OpenArchiver

Shows ingestion sources with error status from your [OpenArchiver](https://github.com/LogicLabs-OU/OpenArchiver) instance.

- `OPENARCHIVER_ADDRESS`
- `INTERNAL_OPENARCHIVER_ADDRESS`
- `OPENARCHIVER_SUPER_API_KEY`
