package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Tgbot *tgbotapi.BotAPI

func StartBot() {
	var err error
	Tgbot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Println(err.Error())
	}
	Tgbot.Debug = true
	log.Printf("Authorized on account %s", Tgbot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := Tgbot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		textMsg := MsgRouter(update)
		msg := tgbotapi.NewMessage(chatID, textMsg)
		Tgbot.Send(msg)
	}

}
