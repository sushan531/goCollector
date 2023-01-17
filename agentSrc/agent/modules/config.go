package modules

import "github.com/spf13/viper"

type WebSocketConfig struct {
	Scheme string
	Host   string
	Port   string
	Path   string
}
type SocketConfig struct {
	Network string
	Host    string
	Port    string
}

type Configuration struct {
	WebSocketConfig
	SocketConfig
	Buffer int
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
