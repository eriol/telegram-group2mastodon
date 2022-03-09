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
   | `MASTODON_TOOT_MAX_CHARACTERS` | Toot max lenght              | 500      |
   | `MASTODON_TOOT_VISIBILITY`     | Default toot visibility      | unlisted |
   | `TELEGRAM_BOT_TOKEN`           | Telegram bot token           | *N/A*    |
   | `TELEGRAM_CHAT_ID`             | Telegram alolowed chat id    | *N/A*    |
   | `TELEGRAM_DEBUG`               | Debug messages from Telegram | False    |

   `MASTODON_TOOT_VISIBILITY` allowed values are: `direct`, `private`, `public`
   and `unlisted`.
   To get `MASTODON_ACCESS_TOKEN` see next point.
4. To get `MASTODON_ACCESS_TOKEN` use the subcommand `authenticate`:
   ```
   telegram-group2mastodon authenticate
   ```
   and follow the istructions.

5. Launch the bot:
   ```
   telegram-group2mastodon run
   ```
