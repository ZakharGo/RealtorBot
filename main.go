package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"log/slog"
	"os"
	config "realtorBot/configs"
	"realtorBot/internal/storage"
	"realtorBot/internal/transport"
)

func main() {

	//init logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//init config
	//err := config.InitConfig()
	//if err != nil {
	//	logger.Error("error set config file", slog.String("error", err.Error()))
	//}
	err := config.InitEnv()
	if err != nil {
		logger.Error("error set env", slog.String("error", err.Error()))
	}

	//New pdb with config data viper.GetString
	pdb, err := storage.NewPostgresDB(storage.Config{
		UserName: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
	})
	if err != nil {
		logger.Error("error init Postgres storage", slog.String("error", err.Error()))
	} else {
		logger.Info("Postgres storage initialized")
	}
	rdb, pong, err := storage.NewRedisStorage(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	if err != nil {
		logger.Error("error init Redis storage", slog.String("error", err.Error()))
	} else {
		logger.Info("Redis storage initialized", slog.String("Ping", pong))
	}
	storage := storage.NewStorage(pdb, rdb)
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
