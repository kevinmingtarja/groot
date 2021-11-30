package main

import (
	"fmt"
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

func (env *Env) sendMessage(chatID int, text string) error {
	// Create new message
	msg := tgbotapi.NewMessage(int64(chatID), text)

	// Send message
	_, err := env.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (env *Env) SendErrorMessage(chatID int, e *ErrorLog) error {
	errorMessage := fmt.Sprintf("Error %d at `%s` endpoint in app `%s`. More details: %s/logs/%d", e.HTTPCode, e.RequestURL, e.AppName, os.Getenv("BASE_URL"), e.ID)

	err := env.sendMessage(chatID, errorMessage)
	if err != nil {
		return err
	}

	return nil
}
