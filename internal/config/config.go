package config

type Config struct {
	SrvAddr     string
	DBAddr      string
	AccrualAddr string
}

var cfg Config

func New() *Config {
	parseFlags()
	parseEnvVars()

	return &cfg
}
