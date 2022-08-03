package config

type Config struct {
	Host          string
	Port          uint16
	Username      string
	Password      string
	RPCPath       string
	Debug         bool
	FileViewerCmd string
}

func New() Config {
	return Config{
		Host:    "localhost",
		Port:    9091,
		RPCPath: "/transmission/rpc",
	}
}
