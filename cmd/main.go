package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"os"
	config "realtorBot/configs"
	"realtorBot/internal/storage"
	"realtorBot/internal/transport"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := config.InitConfig()
	if err != nil {
		logger.Error("error set config file", slog.String("error", err.Error()))
	}
	err = config.InitEnv()
	if err != nil {
		logger.Error("error set env", slog.String("error", err.Error()))
	}

	//New db with config data viper.GetString
	db, err := storage.NewPostgresDB(storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		UserName: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: os.Getenv("passwordDB"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logger.Error("error init storage", slog.String("error", err.Error()))
	}
	storage := storage.NewStorage(db)
	handler := transport.NewHandler(storage)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		logger.Error("error init telegram bot", slog.String("error", err.Error()))
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil || update.CallbackQuery != nil { // If we got a message or callbackQuery
			handler.InitHandler(update, bot, logger)
		}
	}
}
