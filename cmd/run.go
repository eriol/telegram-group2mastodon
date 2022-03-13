package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	mastodonapi "github.com/cking/go-mastodon"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"

	"noa.mornie.org/eriol/telegram-group2mastodon/cfg"
	"noa.mornie.org/eriol/telegram-group2mastodon/mastodon"
	"noa.mornie.org/eriol/telegram-group2mastodon/utils"
)

const (
	MASTODON_ACCESS_TOKEN        = "MASTODON_ACCESS_TOKEN"
	MASTODON_SERVER_ADDRESS      = "MASTODON_SERVER_ADDRESS"
	MASTODON_TOOT_FOOTER         = "MASTODON_TOOT_FOOTER"
	MASTODON_TOOT_MAX_CHARACTERS = "MASTODON_TOOT_MAX_CHARACTERS"
	TELEGRAM_BOT_TOKEN           = "TELEGRAM_BOT_TOKEN"
	TELEGRAM_CHAT_ID             = "TELEGRAM_CHAT_ID"
	TELEGRAM_DEBUG               = "TELEGRAM_DEBUG"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Long: `Start the bot making it connect bot to Telegram and to Mastodon.

Every messages posted in the Telegram groups the bot is in will be posted into
the specified Mastodon account.`,
	Run: func(cmd *cobra.Command, args []string) {
		mastodonInstance := os.Getenv(MASTODON_SERVER_ADDRESS)
		c := mastodonapi.NewClient(&mastodonapi.Config{
			Server:      mastodonInstance,
			AccessToken: os.Getenv(MASTODON_ACCESS_TOKEN),
		})
		log.Println("Crating a new client for mastondon istance:", mastodonInstance)
		allowedTelegramChat := parseTelegramChatID(os.Getenv(TELEGRAM_CHAT_ID))
		log.Println("Allowed telegram chat id:", allowedTelegramChat)

		bot, err := tgbotapi.NewBotAPI(os.Getenv(TELEGRAM_BOT_TOKEN))
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = parseBoolOrFalse(os.Getenv(TELEGRAM_DEBUG))

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 30

		updates := bot.GetUpdatesChan(u)

		for update := range updates {
			chatID := update.Message.Chat.ID
			if chatID != allowedTelegramChat {
				log.Printf("Error: telegram chat %d is not the allowed one: %d\n",
					chatID,
					allowedTelegramChat,
				)
				continue
			}

			if update.Message != nil {
				messageID := update.Message.MessageID
				maxChars := cfg.GetMastodonMaxCharacters()
				tootVisibility := cfg.GetMastodonVisibility()
				tootFooter := os.Getenv(MASTODON_TOOT_FOOTER)

				if update.Message.Text != "" {
					log.Printf("Text message received. Message id: %d\n", messageID)

					text := update.Message.Text
					messages := utils.SplitTextAtChunk(text, maxChars, tootFooter)
					mastodon.PostToots(c, messages, tootVisibility)

				} else if update.Message.Photo != nil {
					log.Printf("Photo received. Message id: %d\n", messageID)

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
						log.Printf("Could not download file: %v", err)
						continue
					}

					mastodon.PostPhoto(
						c,
						file,
						update.Message.Caption,
						maxChars,
						tootVisibility,
					)
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

func parseTelegramChatID(s string) int64 {
	r, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return r
}
