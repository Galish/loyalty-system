package config

import "flag"

func init() {
	flag.StringVar(&cfg.SrvAddr, "a", ":8888", "Server address")
	flag.StringVar(&cfg.DBAddr, "d", "", "DB address")
	flag.StringVar(&cfg.SecretKey, "s", "yvdUuY)HSX}?&b", "JWT signing secret key")
	flag.StringVar(&cfg.AccrualAddr, "r", "", "Accrual system address")
	flag.UintVar(&cfg.AccrualInterval, "i", 1000, "Accrual limiter interval in milliseconds")
	flag.StringVar(&cfg.LogLevel, "l", "info", "Log level")
}

func parseFlags() {
	flag.Parse()
}
