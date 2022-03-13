package cfg

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	MASTODON_ACCESS_TOKEN        = "MASTODON_ACCESS_TOKEN"
	MASTODON_SERVER_ADDRESS      = "MASTODON_SERVER_ADDRESS"
	MASTODON_TOOT_FOOTER         = "MASTODON_TOOT_FOOTER"
	MASTODON_TOOT_MAX_CHARACTERS = "MASTODON_TOOT_MAX_CHARACTERS"
	MASTODON_TOOT_VISIBILITY     = "MASTODON_TOOT_VISIBILITY"
	TELEGRAM_BOT_TOKEN           = "TELEGRAM_BOT_TOKEN"
	TELEGRAM_CHAT_ID             = "TELEGRAM_CHAT_ID"
	TELEGRAM_DEBUG               = "TELEGRAM_DEBUG"
)

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
	return parseMastodonVisibility(os.Getenv(MASTODON_TOOT_VISIBILITY))
}

// Parse Mastodon max characters and return 500 as default in case of errors.
func parseMastodonMaxCharacters(s string) int {
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		return int(n)
	}

	return 500
}

func GetMastodonMaxCharacters() int {
	return parseMastodonMaxCharacters(os.Getenv(MASTODON_TOOT_MAX_CHARACTERS))
}
