package modules

import "github.com/spf13/viper"

type ZmqConfig struct {
	Host    string
	Port    string
	Network string
}

type WebSocketConfig struct {
	Scheme string
	Host   string
	Port   string
	Path   string
}

type Configuration struct {
	ZmqConfig
	WebSocketConfig
	MinHandlers int
}

func ReadConfig(filePath string) Configuration {
	var configuration Configuration

	viper.AddConfigPath(filePath)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	CheckErr(err)
	err = viper.Unmarshal(&configuration)
	CheckErr(err)
	return configuration
}
