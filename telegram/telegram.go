package telegram

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/saeziae/tg2mas-go/utils"
)

func Init(token string) *telego.Bot {
	bot, err := telego.NewBot(token)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		log.Println("Telegram bot started")
	}

	botUser, err := bot.GetMe()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	log.Print("Bot user: ", botUser.ID)

	return bot
}

func Post(bot *telego.Bot, chatId int64, message string) {
	_, err := bot.SendMessage(tu.Message(tu.ID(chatId), message))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Message sent to telegram", chatId)
	}
}

func Listen(bot *telego.Bot, chatID int64, postFuncs ...func(utils.Msg)) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{}, 1)
	updatesParams := telego.GetUpdatesParams{
		AllowedUpdates: []string{"message", "channel_post"},
	}
	updates, _ := bot.UpdatesViaLongPolling(&updatesParams)

	bh, _ := th.NewBotHandler(bot, updates)

	handle := func(bot *telego.Bot, message telego.Message) {
		if chatID == message.Chat.ID {
			log.Print("Message received from Telegram: ", chatID)
			var msg utils.Msg
			if message.Text != "" {
				log.Println("Text message")
				if strings.Contains(message.Text, "!fwdoff") {
					log.Println("There is a forwarding off mark in the message")
					goto no_post
				}
				msg = utils.Msg{Text: message.Text}
			} else if message.MediaGroupID != "" {
				log.Println("Media group message, not supported")
				goto no_post
			} else if message.Photo != nil {
				log.Println("Single photo message")
				photoID := message.Photo[len(message.Photo)-1].FileID
				// Get the photo file
				file, _ := bot.GetFile(&telego.GetFileParams{
					FileID: photoID,
				})
				fileData, err := tu.DownloadFile(bot.FileDownloadURL(file.FilePath))
				if err != nil {
					log.Fatal(err)
					goto no_post
				}
				if strings.Contains(message.Caption, "!fwdoff") {
					log.Println("There is a forwarding off mark in the message")
					goto no_post
				}
				msg = utils.Msg{Text: message.Caption, Media: [][]byte{fileData}}
			} else {
				goto no_post
			}
			for _, postFunc := range postFuncs {
				postFunc(msg)
			}
		no_post:
		}
	}

	bh.HandleChannelPost(handle)
	bh.HandleMessage(handle)

	go func() {
		<-sigs

		log.Println("Stopping...")

		bot.StopLongPolling()
		log.Println("Long polling done")

		bh.Stop()
		log.Println("Bot handler done")

		done <- struct{}{}
	}()

	go bh.Start()
	log.Println("Handling updates...")

	<-done
	log.Println("Done")
}
