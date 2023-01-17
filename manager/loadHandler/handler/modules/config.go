package modules

type Configuration struct {
	ZmqPull
	RedisConn
}

type ZmqPull struct {
	Host     string `default:"127.0.0.1"`
	port     string `default:"5556"`
	Type     string `default:"tcp"`
	Protocol string `default:"PULL"`
}

type RedisConn struct {
	Host     string `default:"127.0.0.1"`
	Port     string `default:"6379"`
	Database int    `default:"0"`
	Password string `default:""`
}
