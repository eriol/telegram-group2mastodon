package cmd

import (
	"log"

	mastodonapi "github.com/cking/go-mastodon"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"

	"noa.mornie.org/eriol/telegram-group2mastodon/cfg"
	"noa.mornie.org/eriol/telegram-group2mastodon/mastodon"
	"noa.mornie.org/eriol/telegram-group2mastodon/utils"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the bot",
	Long: `Start the bot making it connect bot to both Telegram and Mastodon.

Every messages posted in the Telegram groups the bot is in will be posted into
the specified Mastodon account.`,
	Run: func(cmd *cobra.Command, args []string) {
		mastodonInstance := cfg.GetMastodonServerAddress()
		c := mastodonapi.NewClient(&mastodonapi.Config{
			Server:      mastodonInstance,
			AccessToken: cfg.GetMastodonAccessToken(),
		})
		log.Println("Crating a new client for mastondon istance:", mastodonInstance)
		allowedTelegramChat := cfg.GetTelegramChatID()
		log.Println("Allowed telegram chat id:", allowedTelegramChat)

		bot, err := tgbotapi.NewBotAPI(cfg.GetTelegramBotToken())
		if err != nil {
			log.Panic(err)
		}

		bot.Debug = cfg.GetTelegramDebug()

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
				tootFooter := cfg.GetMastodonTootFooter()

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
					file, err := utils.DownloadFile(url)
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
