package mastodon

import (
	"context"
	"log"

	"github.com/mattn/go-mastodon"
	"github.com/saeziae/tg2mas-go/utils"
)

func Init(conf utils.Mastodon) *mastodon.Client {
	config := &mastodon.Config{
		Server:       conf.Server,
		ClientID:     conf.ClientID,
		ClientSecret: conf.CientSecret,
		AccessToken:  conf.AccessToken,
	}

	// Create the client
	c := mastodon.NewClient(config)
	return c
}

func Post(msg utils.Msg, c *mastodon.Client) {
	visibility := "public"
	var mediaIDs []mastodon.ID
	for _, media := range msg.Media {
		upMedia, err := c.UploadMediaFromBytes(context.Background(), media)
		if err != nil {
			log.Fatal(err)
		} else {
			mediaIDs = append(mediaIDs, upMedia.ID)
		}
	}
	toot := mastodon.Toot{
		Status:     msg.Text,
		Visibility: visibility,
		MediaIDs:   mediaIDs,
	}
	_, err := c.PostStatus(context.Background(), &toot)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Post sent to mastodon")
	}
}
