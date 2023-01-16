package modules

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
