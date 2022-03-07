package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cking/go-mastodon"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

const (
	DEBUG                        = "DEBUG"
	TELEGRAM_BOT_TOKEN           = "TELEGRAM_BOT_TOKEN"
	MASTODON_SERVER_ADDRESS      = "MASTODON_SERVER_ADDRESS"
	MASTODON_ACCESS_TOKEN        = "MASTODON_ACCESS_TOKEN"
	MASTODON_TOOT_VISIBILITY     = "MASTODON_TOOT_VISIBILITY"
	MASTODON_TOOT_MAX_CHARACTERS = "MASTODON_TOOT_MAX_CHARACTERS"
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
		max_characters := parseMastodonMaxCharacters(os.Getenv(MASTODON_TOOT_MAX_CHARACTERS))

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

					message := update.Message.Text
					length := len([]rune(message))

					in_reply_to := ""
					for start := 0; start < length; start += max_characters {
						end := start + max_characters
						if end > length {
							end = length
						}

						message_to_post := string([]rune(message)[start:end])
						status, err := c.PostStatus(context.Background(), &mastodon.Toot{
							Status:      message_to_post,
							Visibility:  parseMastodonVisibility(os.Getenv(MASTODON_TOOT_VISIBILITY)),
							InReplyToID: mastodon.ID(in_reply_to),
						})
						if err != nil {
							log.Printf("Could not post status: %v", err)
							continue
						}
						log.Printf("Posted status %s", status.URL)
						in_reply_to = string(status.ID)
					}

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
					log.Printf("Downloading: %s\n", url)
					file, err := downloadFile(url)
					if err != nil {
						log.Printf("Could not post status: %v", err)
						continue
					}
					attachment, err := c.UploadMediaFromReader(
						context.Background(), file)
					if err != nil {
						log.Printf("Could not upload media: %v", err)
						continue
					}
					file.Close()
					log.Printf("Posted attachment %s", attachment.TextURL)

					mediaIds := [...]mastodon.ID{attachment.ID}
					status, err := c.PostStatus(context.Background(), &mastodon.Toot{
						// Write the caption in the toot because it almost probably
						// doesn't describe the image.
						Status:     update.Message.Caption[:max_characters],
						MediaIDs:   mediaIds[:],
						Visibility: parseMastodonVisibility(os.Getenv(MASTODON_TOOT_VISIBILITY)),
					})
					if err != nil {
						log.Printf("Could not post status: %v", err)
						continue
					}
					log.Printf("Posted status %s", status.URL)
				}
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

func downloadFile(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Not able to download %s", url)
	}

	return response.Body, nil
}

// Parse Mastodon max characters and return 500 as default in case of errors.
func parseMastodonMaxCharacters(s string) int {
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		return int(n)
	}

	return 500
}
