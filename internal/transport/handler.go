package transport

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"realtorBot/internal/core"
	"realtorBot/internal/storage"
	"strconv"
	"strings"
	"time"
)

// todo const
type Handler struct {
	storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) InitHandler(tgUpdate tgbotapi.Update, bot *tgbotapi.BotAPI, logger *slog.Logger) {
	if tgUpdate.Message != nil {
		logger.Info("GetMessage:", tgUpdate.Message.Text, "From", tgUpdate.Message.From.UserName)
		message := tgUpdate.Message
		h.mainHandler(core.Message{
			Text:            message.Text,
			UserID:          message.From.ID,
			UserName:        message.From.UserName,
			UserDisplayName: fmt.Sprint(message.From.FirstName + " " + message.From.LastName),
			IsCallback:      false,
			MessageChatID:   message.Chat.ID,
		}, bot, logger)
	} else if tgUpdate.CallbackQuery != nil {
		callback := tgUpdate.CallbackQuery
		logger.Info(fmt.Sprintf("[%s][%v] Callback: %s", tgUpdate.CallbackQuery.From.UserName, tgUpdate.CallbackQuery.From.ID, tgUpdate.CallbackQuery.Data))
		h.mainHandler(core.Message{
			UserID:          callback.From.ID,
			UserName:        callback.From.UserName,
			UserDisplayName: fmt.Sprint(callback.From.FirstName + " " + callback.From.LastName),
			IsCallback:      true,
			CallbackMsgID:   callback.ID,
			CallbackData:    callback.Data,
			CallbackChatID:  callback.From.ID,
		}, bot, logger)

	}
}

// todo var
func (h *Handler) mainHandler(msg core.Message, bot *tgbotapi.BotAPI, logger *slog.Logger) {
	if msg.IsCallback == true {
		switch msg.CallbackData {
		case "newflat":
			if err := h.storage.Cache.Create(msg.CallbackData); err != nil {
				logger.Error("error:", slog.String("error in flat callbackQuery", err.Error()))
				core.SendMessageTg(msg.CallbackChatID, core.ErrorAnswer, bot)
			} else {
				core.SendMessageTg(msg.CallbackChatID, core.NewFlatAnswerCallback, bot)
			}
		case "getallflat":
			flats, err := h.storage.Flat.GetAll()
			if err != nil {
				logger.Error("error:", slog.String("error in get all Flats", err.Error()))
				core.SendMessageTg(msg.CallbackChatID, core.ErrorAnswer, bot)
			} else {
				for i := 0; i < len(flats); i++ {
					core.SendMessageTg(msg.CallbackChatID, flats[i], bot)
				}
			}
		case "newcount":
			if err := h.storage.Cache.Create(msg.CallbackData); err != nil {
				logger.Error("error:", slog.String("error in new count callbackQuery", err.Error()))
			} else {
				core.SendMessageTg(msg.CallbackChatID, core.NewCountAnswerCallback, bot)
			}

		}
	} else {
		switch msg.Text {
		//todo command /start try
		case "/start":
			if err := core.NewStartInlineBtn(msg.MessageChatID, bot); err != nil {
				logger.Error("error send inline Button", slog.String("error", err.Error()))
				core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
			}
		default:
			if err := h.storage.Delete(""); err != nil {
				logger.Error("error delete cache", slog.String("error", err.Error()))
			}
			data, err := h.storage.Cache.Get()
			if err != nil {
				logger.Error("error get cache", slog.String("error", err.Error()))
				core.SendMessageTg(msg.MessageChatID, core.NotFoundCommand, bot)
			} else {
				h.HandleInputData(msg, bot, data, logger)
			}
		}
	}
}

func (h *Handler) HandleInputData(msg core.Message, bot *tgbotapi.BotAPI, data string, logger *slog.Logger) {
	switch data {
	case "newflat":
		if err := h.storage.Cache.Delete(data); err != nil {
			logger.Error("error delete cache", slog.String("error", err.Error()))
		}
		err := h.storage.Flat.Create(msg.Text)
		if err != nil {
			logger.Error("error create flat", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.RepeatingMeaning, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.TaskCompletedSuccessfully, bot)
		}

	case "newcount":
		if err := h.storage.Cache.Delete(data); err != nil {
			logger.Error("error delete cache", slog.String("error", err.Error()))
		}
		txt := strings.Split(msg.Text, " ")
		if len(txt) != 2 {
			core.SendMessageTg(msg.MessageChatID, core.ErrorInputData, bot)
			if err := h.storage.Cache.Delete("newflat"); err != nil {
				logger.Error("error delete cache", slog.String("error", err.Error()))
			}
			break
		}
		numb := txt[0]
		count, err := strconv.Atoi(txt[1])
		if err != nil {
			core.SendMessageTg(msg.MessageChatID, core.ErrorInputData, bot)
			if err := h.storage.Cache.Delete(data); err != nil {
				logger.Error("error delete cache", slog.String("error", err.Error()))
			}
			break
		}
		date := time.Now()
		err = h.storage.Count.Create(numb, count, date)
		if err != nil {
			logger.Error("error create count", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.TaskCompletedSuccessfully, bot)
		}

	}
}
