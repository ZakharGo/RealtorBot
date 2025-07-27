package transport

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
	"os/signal"
	"realtorBot/internal/core"
	"realtorBot/internal/storage"
	"strconv"
	"strings"
	"syscall"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	if tgUpdate.Message != nil {
		logger.Info(fmt.Sprintf("[%s][%v] Message: %s", tgUpdate.Message.From.UserName, tgUpdate.Message.From.ID, tgUpdate.Message.Text))
		message := tgUpdate.Message
		h.mainHandler(ctx, core.Message{
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
		h.mainHandler(ctx, core.Message{
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
func (h *Handler) mainHandler(ctx context.Context, msg core.Message, bot *tgbotapi.BotAPI, logger *slog.Logger) {
	//init Context
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	data := msg.Text
	userID := fmt.Sprint(msg.UserID)

	switch msg.Text {
	case "NewFlat":
		if err := h.storage.Cache.Create(ctx, data, userID); err != nil {
			logger.Error("error:", slog.String("error in flat callbackQuery", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.NumberFlatAnswerCallback, bot)
		}
	case "DeleteFlat":
		if err := h.storage.Cache.Create(ctx, data, userID); err != nil {
			logger.Error("error:", slog.String("error in flat callbackQuery", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.NumberFlatAnswerCallback, bot)
		}
	case "GetAllFlat":
		getAllString, err := h.storage.Flat.GetAll()
		if err != nil {
			logger.Error("error:", slog.String("error in get all Flats", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		} else {
			for _, v := range getAllString {
				core.SendMessageTg(msg.MessageChatID, v, bot)
			}
		}
	case "NewCount":
		if err := h.storage.Cache.Create(ctx, data, userID); err != nil {
			logger.Error("error:", slog.String("error in new count callbackQuery", err.Error()))
		} else {
			core.SendMessageTg(msg.MessageChatID, core.NewCountAnswerCallback, bot)
		}
	case "GetLastCount":
		if err := h.storage.Cache.Create(ctx, data, userID); err != nil {
			logger.Error("error:", slog.String("error in get last count callbackQuery", err.Error()))
		} else {
			core.SendMessageTg(msg.MessageChatID, core.NewGetAmountOfPayment, bot)
		}

	case "DeleteLastCount":
		flat, count, err := h.storage.Count.DeleteLastCount()
		if err != nil {
			logger.Error("error:", slog.String("error in delete last count callbackQuery", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, fmt.Sprintf("Запись %v квартиры %s удалена", count, flat), bot)
		}
	case "Amount of payment":
		if err := h.storage.Cache.Create(ctx, data, userID); err != nil {
			logger.Error("error:", slog.String("error in Amount of payment callbackQuery", err.Error()))
		} else {
			core.SendMessageTg(msg.MessageChatID, core.NewGetAmountOfPayment, bot)
		}
	case "/start":
		if err := core.Mainkeyboard(msg.MessageChatID, bot); err != nil {
			logger.Error("error send inline Button", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		}
	default:

		data, err := h.storage.Cache.Get(ctx, userID)
		if err != nil {
			logger.Error("error get cache", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.NotFoundCommand, bot)
		} else {
			h.HandleInputData(ctx, msg, bot, data, logger)
		}
	}
}

func (h *Handler) HandleInputData(ctx context.Context, msg core.Message, bot *tgbotapi.BotAPI, data string, logger *slog.Logger) {
	switch data {
	case "NewFlat":

		err := h.storage.Flat.Create(msg.Text)
		if err != nil {
			logger.Error("error create flat", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.RepeatingMeaning, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.TaskCompletedSuccessfully, bot)
		}

	case "DeleteFlat":

		if err := h.storage.Flat.Delete(msg.Text); err != nil {
			logger.Error("error delete flat", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.FlatNotFound, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.TaskCompletedSuccessfully, bot)
		}

	case "NewCount":

		txt := strings.Split(msg.Text, " ")
		if len(txt) != 2 {
			core.SendMessageTg(msg.MessageChatID, core.ErrorInputData, bot)

			break
		}
		numb := txt[0]
		count, err := strconv.Atoi(txt[1])
		if err != nil {
			core.SendMessageTg(msg.MessageChatID, core.ErrorInputData, bot)
			break
		}
		date := time.Now()
		err = h.storage.Count.Create(numb, count, date)
		if err != nil {
			logger.Error("error create count", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.FlatNotFound, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, core.TaskCompletedSuccessfully, bot)
		}
	case "GetLastCount":
		count, err := h.storage.Count.GetLast(msg.Text)
		if err != nil {
			logger.Error("error get last count", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.FlatNotFound, bot)
		} else {
			core.SendMessageTg(msg.MessageChatID, fmt.Sprintf("Последня запись квартиры %s = %v", msg.Text, count), bot)
		}
	case "Amount of payment":
		numb := msg.Text
		LastCount, err := h.storage.Count.GetLast(numb)
		if err != nil {
			logger.Error("error get last count", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		}
		PenultCount, err := h.storage.Count.GetPenult(numb)
		if err != nil {
			logger.Error("error get penult count", slog.String("error", err.Error()))
			core.SendMessageTg(msg.MessageChatID, core.ErrorAnswer, bot)
		}
		amount := float64(LastCount-PenultCount) * core.PriceOfElectricity
		if amount > 0.0 {
			core.SendMessageTg(msg.MessageChatID, fmt.Sprintf("Здравствуйте, показаниe счетчика %v к оплате %v рублей", LastCount, amount), bot)
		}

	}

}
