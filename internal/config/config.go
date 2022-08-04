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

const (
	ProgramName string = "gotem"
)

func New() Config {
	return Config{
		Host:    "localhost",
		Port:    9091,
		RPCPath: "/transmission/rpc",
	}
}
