package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func SendMessageTg(chatId int64, text string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, text)
	bot.Send(msg)
}
