package core

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log/slog"
)

// nolint
func NewStartInlineBtn(ChatId int64, bot *tgbotapi.BotAPI) error {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("NewFlat", "newflat"),
			tgbotapi.NewInlineKeyboardButtonData("GetAllFlat", "getallflat"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("NewCount", "newcount")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Amount of Payment", "amountofpayment")),
	)
	msg := tgbotapi.NewMessage(ChatId, "Нажми на кнопку:")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		return err
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
