package modules

type RedisConfig struct {
	Host string
	Port string
}

type ZmqConfig struct {
	ConnType string
	Host     string
	Port     string
}

type RedisStream struct {
	Stream string
	Dummy  bool
}

type Configuration struct {
	RedisConfig
	ZmqConfig
	RedisStream
}
