package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func registerSubscriber(firstName string, lastName string, chatID int64) string {

	subscriber := Subscriber{Name: firstName + "_" + lastName, ChatID: chatID, FirstName: firstName, LastName: lastName}
	log.Println(subscriber)
	return "Subscriber registered"
}

func subscribe(firstName string, lastName string, chatID int64, serviceName string) string {
	// check if service exists
	service, err := GetService(serviceName)
	if err != nil {
		return "Service not found"
	}
	// check if subscriber exists
	subscriber, err := GetSubscriber(chatID)
	if err != nil {
		return "Subscriber not found"
	}
	s := Subscriptions{SubscriberID: subscriber.ID, ServiceID: service.ID}
	err = s.Create()
	if err != nil {
		return "Subscription failed because " + err.Error()
	}
	// check if subscriber is already subscribed to service
	// create subscription
	return "Subscription successful"
}

// commands -> register, subscribe, help,
func MsgRouter(update tgbotapi.Update) string {
	msg := update.Message.Text

	words := strings.Split(msg, " ")
	switch words[0] {
	case "/register":
		return registerSubscriber(update.Message.Chat.FirstName, update.Message.Chat.LastName, update.Message.Chat.ID)
	case "/subscribe":
		if len(words) < 2 {
			return "Please specify a service"
		}
		return subscribe(update.Message.Chat.FirstName, update.Message.Chat.LastName, update.Message.Chat.ID, words[1])
	case "/help":
		return "Help"
	default:
		return "Unknown command"
	}
}
