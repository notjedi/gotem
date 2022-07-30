package config

type Config struct {
	Host          string
	Port          uint16
	Username      string
	Password      string
	RpcPath       string
	FileViewerCmd string
}

func New() Config {
	return Config{
		Host: "localhost",
		Port: 9091,
	}
}
