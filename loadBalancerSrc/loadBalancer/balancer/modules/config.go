package modules

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
