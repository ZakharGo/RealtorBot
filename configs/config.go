package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
