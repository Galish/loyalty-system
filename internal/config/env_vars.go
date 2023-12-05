package config

import "os"

func parseEnvVars() {
	if runAddr := os.Getenv("RUN_ADDRESS"); runAddr != "" {
		cfg.SrvAddr = runAddr
	}

	if dbAddr := os.Getenv("DATABASE_URI"); dbAddr != "" {
		cfg.DBAddr = dbAddr
	}

	if accrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualAddr != "" {
		cfg.AccrualAddr = accrualAddr
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
}
