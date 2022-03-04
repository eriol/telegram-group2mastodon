package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cking/go-mastodon"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

const (
	DEBUG                    = "DEBUG"
	TELEGRAM_BOT_TOKEN       = "TELEGRAM_BOT_TOKEN"
	MASTODON_SERVER_ADDRESS  = "MASTODON_SERVER_ADDRESS"
	MASTODON_ACCESS_TOKEN    = "MASTODON_ACCESS_TOKEN"
	MASTODON_TOOT_VISIBILITY = "MASTODON_TOOT_VISIBILITY"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Long: `Start the bot making it connect bot to Telegram and to Mastodon.

Every messages posted in the Telegram groups the bot is in will be posted into
the specified Mastodon account.`,
	Run: func(cmd *cobra.Command, args []string) {
		mastodon_instance := os.Getenv(MASTODON_SERVER_ADDRESS)
		c := mastodon.NewClient(&mastodon.Config{
			Server:      mastodon_instance,
			AccessToken: os.Getenv(MASTODON_ACCESS_TOKEN),
		})
		log.Println("Crating a new client for mastondon istance:", mastodon_instance)

		bot, err := tgbotapi.NewBotAPI(os.Getenv(TELEGRAM_BOT_TOKEN))
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = parseBoolOrFalse(os.Getenv(DEBUG))

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 30

		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			if update.Message != nil {

				if update.Message.Text != "" {
					log.Println("Text message received.")
					status, err := c.PostStatus(context.Background(), &mastodon.Toot{
						Status:     update.Message.Text,
						Visibility: parseMastodonVisibility(os.Getenv(MASTODON_TOOT_VISIBILITY)),
					})

					if err != nil {
						log.Fatalf("Could not post status: %v", err)
					}

					log.Printf("Posted status %s", status.URL)
				} else if update.Message.Photo != nil {
					log.Println("Photo received.")
					// Telegram provides multiple sizes of photo, just take the
					// biggest.
					biggest_photo := tgbotapi.PhotoSize{FileSize: 0}
					for _, photo := range update.Message.Photo {

						if photo.FileSize > biggest_photo.FileSize {
							biggest_photo = photo
						}

					}
					url, _ := bot.GetFileDirectURL(biggest_photo.FileID)
					fmt.Println(url)

				}
				// fmt.Printf("%#v\n\n", update.Message)

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func parseBoolOrFalse(s string) bool {
	r, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return r
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
