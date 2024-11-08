package main

import (
	"fmt"

	"github.com/saeziae/tg2mas-go/mastodon"
	"github.com/saeziae/tg2mas-go/telegram"
	"github.com/saeziae/tg2mas-go/utils"
)

const (
	version = "0.0.1"
)

func main() {
	fmt.Println("tg2mas", "Ver", version)
	utils.PrintPreamable() // Print the license
	utils.LoadConfig()
	telegramBot := telegram.Init(utils.Config.Telegram.Token)
	mastodonBot := mastodon.Init(utils.Config.Mastodon)
	//wrap the post function for passing as a parameter
	mastodonPost := func(msg utils.Msg) {
		mastodon.Post(msg, mastodonBot)
	}
	telegram.Listen(telegramBot, utils.Config.Telegram.ChatID, mastodonPost)

}
