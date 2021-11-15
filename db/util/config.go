package util

import "github.com/spf13/viper"

// Config stores all configuration of the application
// Viper read configration from a config file or enviroment variable
type Config struct {
	DBDriven      string `mapstructure:"DB_DRIVEN"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}