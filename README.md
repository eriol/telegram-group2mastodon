# telegram-group2mastodon

telegram-group2mastodon is a Telegram bot to post messages from a Telegram
group to Mastodon. It's written in Go and it is relased under the AGPL3+.
Since the bot must be able to read all the messages you should use it only on
public groups.

## Installation

1. Build the bot (for example if version is v0.1.0 and you are on a unix system):
   ```
   cd /tmp; GOPATH=/tmp/go go install noa.mornie.org/eriol/telegram-group2mastodon@v0.1.0
   ```
   You will find the binary at `/tmp/go/bin/telegram-group2mastodon`. You don't
   need to set `GOPATH` if you already have it set and you are fine having the
   bot installed there.

2. Create your bot using Telegram @BotFather.

3. The bot uses environment variables as configuration, you have to export
   the following (except variables with a default) before starting it:
   | Variable                       | Meaning                      | Default  |
   |--------------------------------|------------------------------|----------|
   | `MASTODON_ACCESS_TOKEN`        | Mastodon access token        | *N/A*    |
   | `MASTODON_SERVER_ADDRESS`      | Mastodon server address      | *N/A*    |
   | `MASTODON_TOOT_FOOTER   `      | Footer of each toot          | ""       |
   | `MASTODON_TOOT_MAX_CHARACTERS` | Toot max lenght              | 500      |
   | `MASTODON_TOOT_VISIBILITY`     | Default toot visibility      | unlisted |
   | `TELEGRAM_BOT_TOKEN`           | Telegram bot token           | *N/A*    |
   | `TELEGRAM_CHAT_ID`             | Telegram alolowed chat id    | *N/A*    |
   | `TELEGRAM_DEBUG`               | Debug messages from Telegram | False    |

   `MASTODON_TOOT_VISIBILITY` allowed values are: `direct`, `private`, `public`
   and `unlisted`.
   `MASTODON_TOOT_FOOTER` default is an empty string that disable the footer
   feature, no character are automatically added so you may want to add a space
   before your footer (or a new line).
   To get `MASTODON_ACCESS_TOKEN` see next point.
4. To get `MASTODON_ACCESS_TOKEN` use the subcommand `authenticate`:
   ```
   telegram-group2mastodon authenticate <https://your.instance>
   ```
   and follow the istructions.

5. Launch the bot:
   ```
   telegram-group2mastodon run
   ```

## Deploy using Docker/Podman

Images are automatically built after each push on quay.io, here the list of
available tags: https://quay.io/repository/eriol/telegram-group2mastodon?tab=tags

### docker

*Suppose you want to use the latest tag.*

```
docker run \
    --env MASTODON_ACCESS_TOKEN=$MASTODON_ACCESS_TOKEN \
    --env MASTODON_SERVER_ADDRESS=$MASTODON_SERVER_ADDRESS \
    --env TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN \
    --env TELEGRAM_CHAT_ID=$TELEGRAM_CHAT_ID \
    quay.io/eriol/telegram-group2mastodon:latest
```

### podman

*Suppose you want to use the latest tag.*

```
podman run \
    --env MASTODON_ACCESS_TOKEN=$MASTODON_ACCESS_TOKEN \
    --env MASTODON_SERVER_ADDRESS=$MASTODON_SERVER_ADDRESS \
    --env TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN \
    --env TELEGRAM_CHAT_ID=$TELEGRAM_CHAT_ID \
    quay.io/eriol/telegram-group2mastodon:latest
```

### docker-compose

```
---
version: "3.9"

services:
  bot:
    image: quay.io/eriol/telegram-group2mastodon:latest
    environment:
      - MASTODON_ACCESS_TOKEN=${MASTODON_ACCESS_TOKEN}
      - MASTODON_SERVER_ADDRESS=${MASTODON_SERVER_ADDRESS}
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - TELEGRAM_CHAT_ID=${TELEGRAM_CHAT_ID}
```

## Deploy on metal

You have to build the bot and use a init system of your choice to handle it.

### systemd

Suppose you put the bot in `/srv/tg2m/telegram-group2mastodon` an example unit
file using an user `tg2mbot` could be:
```
[Unit]
Description=a Telegram bot to post messages from a Telegram group to Mastodon
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/srv/tg2m/telegram-group2mastodon run
User=tg2mbot
Group=tg2mbot
Restart=on-failure
Environment=MASTODON_ACCESS_TOKEN=<MASTODON_ACCESS_TOKEN>
Environment=MASTODON_SERVER_ADDRESS=<MASTODON_SERVER_ADDRESS>
Environment=TELEGRAM_BOT_TOKEN=<TELEGRAM_BOT_TOKEN>
Environment=TELEGRAM_CHAT_ID=<TELEGRAM_CHAT_ID>

[Install]
WantedBy=multi-user.target
```

where you have to put the real values for `<MASTODON_ACCESS_TOKEN>`,
`<MASTODON_SERVER_ADDRESS>`, `<TELEGRAM_BOT_TOKEN>` and `<TELEGRAM_CHAT_ID>`.
