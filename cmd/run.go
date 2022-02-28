/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/cking/go-mastodon"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/cobra"
)

const (
	DEBUG                   = "DEBUG"
	TELEGRAM_BOT_TOKEN      = "TELEGRAM_BOT_TOKEN"
	MASTODON_SERVER_ADDRESS = "MASTODON_SERVER_ADDRESS"
	MASTODON_CLIENT_ID      = "MASTODON_CLIENT_ID"
	MASTODON_SECRET         = "MASTODON_SECRET"
	MASTODON_ACCESS_TOKEN   = "MASTODON_ACCESS_TOKEN"
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
			Server:       mastodon_instance,
			ClientID:     os.Getenv(MASTODON_CLIENT_ID),
			ClientSecret: os.Getenv(MASTODON_SECRET),
		})
		log.Println("Crating a new client for mastondon istance:", mastodon_instance)

		err := c.AuthenticateToken(context.Background(), os.Getenv(MASTODON_ACCESS_TOKEN), os.Getenv(MASTODON_REDIRECT_URI))
		if err != nil {
			log.Fatal(err)
		}

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
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				log.Println(update.Message)

				status, err := c.PostStatus(context.Background(), &mastodon.Toot{
					Status:     update.Message.Text,
					Visibility: "unlisted",
				})

				if err != nil {
					log.Fatalf("Could not post status: %v", err)
				}

				log.Printf("Posted status %s", status.URL)
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
