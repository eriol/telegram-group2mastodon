package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cking/go-mastodon"
	"github.com/spf13/cobra"
)

const MASTODON_REDIRECT_URI = "urn:ietf:wg:oauth:2.0:oob"

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Authenticate the bot against Mastodon",
	Long:  "Authenticate against the specified Mastodon instance.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mastodon_instance_address := args[0]

		app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
			Server:     mastodon_instance_address,
			ClientName: "tg2mastodon",
			Scopes:     "read write follow",
			Website:    "https://noa.mornie.org/eriol/tg2mastodon",
		})

		if err != nil {
			log.Fatal(err)
		}

		client := mastodon.NewClient(&mastodon.Config{
			Server:       mastodon_instance_address,
			ClientID:     app.ClientID,
			ClientSecret: app.ClientSecret,
		})

		fmt.Printf("Please open the following URL to authenticate: %s\n",
			app.AuthURI)
		fmt.Printf("And paste the authentication code here and press enter: ")

		reader := bufio.NewReader(os.Stdin)
		authCode, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		authCode = strings.TrimSuffix(authCode, "\n")

		err = client.AuthenticateToken(
			context.Background(),
			authCode,
			MASTODON_REDIRECT_URI)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Authentication success!")
		fmt.Println("Please exporte the following environment variables:")
		fmt.Printf("export MASTODON_SERVER_ADDRESS='%s'\n", mastodon_instance_address)
		fmt.Printf("export MASTODON_ACCESS_TOKEN='%s'\n", client.Config.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(authenticateCmd)
}
