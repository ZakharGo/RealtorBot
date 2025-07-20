package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log/slog"
)

// nolint
func Mainkeyboard(ChatId int64, bot *tgbotapi.BotAPI) error {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("NewFlat"),
			tgbotapi.NewKeyboardButton("DeleteFlat"),
			tgbotapi.NewKeyboardButton("GetAllFlat"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("NewCount"),
			tgbotapi.NewKeyboardButton("GetLastCount"),
			tgbotapi.NewKeyboardButton("DeleteLastCount"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Amount of payment"),
		),
	)
	msg := tgbotapi.NewMessage(ChatId, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func AllflatBtn(flats []string, ChatId int64, bot *tgbotapi.BotAPI) error {
	for _, flat := range flats {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(flat, flat),
			),
		)
		msg := tgbotapi.NewMessage(ChatId, "")
		msg.ReplyMarkup = keyboard
		_, err := bot.Send(msg)
		if err != nil {
			return err
			break
		}
	}
	return nil
}

// todo for main keybord try later
func DeleteInlineBtn(userID int64, msgID int, sourceText string, bot *tgbotapi.BotAPI, logger *slog.Logger) error {
	msg := tgbotapi.NewEditMessageText(userID, msgID, sourceText)
	_, err := bot.Send(msg)
	if err != nil {
		logger.Error("Ошибка отправки сообщения", "err", err)
		return errors.Wrap(err, "client.Send remove inline-buttons")
	}
	return nil
}
