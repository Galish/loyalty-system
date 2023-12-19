package config

import "flag"

func init() {
	flag.StringVar(&cfg.SrvAddr, "a", ":8888", "Server address")
	flag.StringVar(&cfg.DBAddr, "d", "", "DB address")
	flag.StringVar(&cfg.SecretKey, "s", "yvdUuY)HSX}?&b", "JWT signing secret key")
	flag.StringVar(&cfg.AccrualAddr, "r", "", "Accrual system address")
	flag.StringVar(&cfg.LogLevel, "l", "debug", "Log level")
}

func parseFlags() {
	flag.Parse()
}
