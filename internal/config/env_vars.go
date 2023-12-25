package config

import (
	"os"
	"strconv"
)

func parseEnvVars() {
	if runAddr := os.Getenv("RUN_ADDRESS"); runAddr != "" {
		cfg.SrvAddr = runAddr
	}

	if dbAddr := os.Getenv("DATABASE_URI"); dbAddr != "" {
		cfg.DBAddr = dbAddr
	}

	if secretKey := os.Getenv("SECRET_KEY"); secretKey != "" {
		cfg.SecretKey = secretKey
	}

	if accrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualAddr != "" {
		cfg.AccrualAddr = accrualAddr
	}

	if accrualInterval := os.Getenv("ACCRUAL_LIMITER_INTERVAL"); accrualInterval != "" {
		interval, _ := strconv.Atoi(accrualInterval)
		cfg.AccrualInterval = uint(interval)
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}
