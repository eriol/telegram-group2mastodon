package cfg

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	mastodonAccessToken       = "MASTODON_ACCESS_TOKEN"
	mastodonServerAddress     = "MASTODON_SERVER_ADDRESS"
	mastodonTootFooter        = "MASTODON_TOOT_FOOTER"
	mastodonTootMaxCharacters = "MASTODON_TOOT_MAX_CHARACTERS"
	mastodonTootVisibility    = "MASTODON_TOOT_VISIBILITY"
	telegramBotToken          = "TELEGRAM_BOT_TOKEN"
	telegramChatID            = "TELEGRAM_CHAT_ID"
	telegramDebug             = "TELEGRAM_DEBUG"
)

func GetMastodonAccessToken() string {
	return os.Getenv(mastodonAccessToken)
}

func GetMastodonServerAddress() string {
	return os.Getenv(mastodonServerAddress)
}

func GetMastodonTootFooter() string {
	return os.Getenv(mastodonTootFooter)
}

// Parse Mastodon max characters and return 500 as default in case of errors.
func parseMastodonMaxCharacters(s string) int {
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		return int(n)
	}

	return 500
}

// Return the max characters allowed for toots.
func GetMastodonMaxCharacters() int {
	return parseMastodonMaxCharacters(os.Getenv(mastodonTootMaxCharacters))
}

// Check the specified Mastodon visibility and return it if valid or return
// unlisted if it's not valid.
// The specified string will be cheched case unsensitive.
func parseMastodonVisibility(s string) string {
	s = strings.ToLower(s)
	// Keep sorted since we search inside.
	visibilities := []string{"direct", "private", "public", "unlisted"}
	r := sort.SearchStrings(visibilities, s)
	if r < len(visibilities) && visibilities[r] == s {
		return s
	}

	return "unlisted"
}

// Return configured Mastodon visibility for toot.
func GetMastodonVisibility() string {
	return parseMastodonVisibility(os.Getenv(mastodonTootVisibility))
}

func parseBoolOrFalse(s string) bool {
	r, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return r
}

// Return if Telegram bot will run in debug mode or not.
func GetTelegramDebug() bool {
	return parseBoolOrFalse(os.Getenv(telegramDebug))
}

func GetTelegramBotToken() string {
	return os.Getenv(telegramBotToken)
}

// Parse the telegram chat or return 0.
func parseTelegramChatID(s string) int64 {
	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return r
}

// Return the configured Telegram chat ID.
func GetTelegramChatID() int64 {
	return parseTelegramChatID(os.Getenv(telegramChatID))
}
