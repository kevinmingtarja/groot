package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

func InitializeBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return bot, nil
}

func (env *Env) SendMessage(chatID int, text string) error {
	// Create new message
	msg := tgbotapi.NewMessage(int64(chatID), text)

	// Send message
	_, err := env.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
