package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// todo no init() if possible
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	// Enable viper to read env variables
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Error("config read failed")
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringOrDefault(key string, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}
