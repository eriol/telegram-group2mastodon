package mastodon

import (
	"context"
	"io"
	"log"

	mastodonapi "github.com/cking/go-mastodon"
)

// Post one or more toots.
func PostToots(
	client *mastodonapi.Client,
	messages []string,
	visibility string,
	inReplyTo string) {
	for _, message := range messages {
		status, err := client.PostStatus(context.Background(), &mastodonapi.Toot{
			Status:      message,
			Visibility:  visibility,
			InReplyToID: mastodonapi.ID(inReplyTo),
		})
		if err != nil {
			log.Printf("Could not post status: %v", err)
			continue
		}
		log.Printf("Posted status %s", status.URL)
		inReplyTo = string(status.ID)
	}
}

// Post a photo on mastodon using the caption as text of the toot.
// If the caption is bigger than the max toot size, we post multiple toots.
func PostPhoto(
	client *mastodonapi.Client,
	file io.ReadCloser,
	messages []string,
	visibility string) {
	attachment, err := client.UploadMediaFromReader(
		context.Background(), file)
	if err != nil {
		log.Printf("Could not upload media: %v", err)
	}
	file.Close()
	log.Printf("Posted attachment %s", attachment.TextURL)

	mediaIds := [...]mastodonapi.ID{attachment.ID}
	// 1. Post the photo with the first part of the caption.
	status, err := client.PostStatus(context.Background(), &mastodonapi.Toot{
		Status:     messages[0],
		MediaIDs:   mediaIds[:],
		Visibility: visibility,
	})
	if err != nil {
		log.Printf("Could not post status: %v", err)
	}
	log.Printf("Posted status %s", status.URL)
	// 2. Post the remaining caption if it exists.
	PostToots(client, messages[1:], visibility, string(status.ID))
}
