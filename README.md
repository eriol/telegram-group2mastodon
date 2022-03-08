# telegram-group2mastodon

telegram-group2mastodon is a Telegram bot to post messages from a Telegram
group to Mastodon. You should only use on public groups and set the bot to be
able to read all the messages.

**This is still a Work In Progress, please don't use at the moment, since also
documentation is missing.**

## TODO

- [x] Publish text messages from Telegram to Mastodon.
- [x] Handle text > 500 characters (or mastodon instance limit).
- [x] Handle messages with 1 photo and caption.
- [x] Handle messages with multiple photos.
- [x] Add configuration variable to set mastodon message visibility.
- [x] Set the Telegram channel allowed to use the bot: we don't want spam from
      someone that just add the bot somewhere.
