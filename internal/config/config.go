package config

type Config struct {
	Host          string
	Port          uint16
	Username      string
	Password      string
	RpcPath       string
	Debug         bool
	FileViewerCmd string
}

func New() Config {
	return Config{
		Host:    "localhost",
		Port:    9091,
		RpcPath: "/transmission/rpc",
	}
}
