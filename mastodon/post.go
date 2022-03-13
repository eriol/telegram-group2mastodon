package mastodon

import (
	"context"
	"io"
	"log"

	mastodonapi "github.com/cking/go-mastodon"
)

// Post one or more toots.
func PostToots(client *mastodonapi.Client, messages []string, visibility string) {
	in_reply_to := ""
	for _, message := range messages {
		status, err := client.PostStatus(context.Background(), &mastodonapi.Toot{
			Status:      message,
			Visibility:  visibility,
			InReplyToID: mastodonapi.ID(in_reply_to),
		})
		if err != nil {
			log.Printf("Could not post status: %v", err)
			continue
		}
		log.Printf("Posted status %s", status.URL)
		in_reply_to = string(status.ID)
	}
}

// Post a photo on mastodon with caption.
func PostPhoto(
	client *mastodonapi.Client,
	file io.ReadCloser,
	caption string,
	maxCharacters int,
	visibility string) {
	attachment, err := client.UploadMediaFromReader(
		context.Background(), file)
	if err != nil {
		log.Printf("Could not upload media: %v", err)
	}
	file.Close()
	log.Printf("Posted attachment %s", attachment.TextURL)

	mediaIds := [...]mastodonapi.ID{attachment.ID}
	if len(caption) > maxCharacters {
		caption = caption[:maxCharacters]
	}
	status, err := client.PostStatus(context.Background(), &mastodonapi.Toot{
		// Write the caption in the toot because it almost probably
		// doesn't describe the image.
		Status:     caption,
		MediaIDs:   mediaIds[:],
		Visibility: visibility,
	})
	if err != nil {
		log.Printf("Could not post status: %v", err)
	}
	log.Printf("Posted status %s", status.URL)
}
