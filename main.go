package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
	err := config.InitConfig()
	if err != nil {
		logger.Error("error set config file", slog.String("error", err.Error()))
	}
	err = config.InitEnv()
	if err != nil {
		logger.Error("error set env", slog.String("error", err.Error()))
	}

	//New pdb with config data viper.GetString
	pdb, err := storage.NewPostgresDB(storage.Config{
		Host:     viper.GetString("pdb.host"),
		Port:     viper.GetString("pdb.port"),
		UserName: viper.GetString("pdb.username"),
		DBName:   viper.GetString("pdb.dbname"),
		Password: os.Getenv("passwordDB"),
		SSLMode:  viper.GetString("pdb.sslmode"),
	})
	if err != nil {
		logger.Error("error init Postgres storage", slog.String("error", err.Error()))
	} else {
		logger.Info("Postgres storage initialized")
	}
	rdb, pong, err := storage.NewRedisStorage(&redis.Options{
		Addr:     viper.GetString("rdb.host"),
		Password: viper.GetString("rdb.password"),
		DB:       viper.GetInt("rdb.db"),
	})
	if err != nil {
		logger.Error("error init Redis storage", slog.String("error", err.Error()))
	}
	logger.Info("Redis storage initialized", slog.String("Ping", pong))

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
