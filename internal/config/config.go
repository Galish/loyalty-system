package config

type Config struct {
	SrvAddr         string
	DBAddr          string
	SecretKey       string
	AccrualAddr     string
	AccrualInterval uint
	LogLevel        string
}

var cfg Config

func New() *Config {
	parseFlags()
	parseEnvVars()

	return &cfg
}
