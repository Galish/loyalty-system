package config

import "flag"

func init() {
	flag.StringVar(&cfg.SrvAddr, "a", ":8888", "Server address")
	flag.StringVar(&cfg.DBAddr, "d", "", "DB address")
	flag.StringVar(&cfg.AccrualAddr, "r", "", "Accrual system address")
}

func parseFlags() {
	flag.Parse()
}
