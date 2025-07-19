package config

import (
	"github.com/joho/godotenv"
)

//	func InitConfig() error {
//		viper.AddConfigPath("configs")
//		viper.SetConfigName("config")
//		return viper.ReadInConfig()
//	}
func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
