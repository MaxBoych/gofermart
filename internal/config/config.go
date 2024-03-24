package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr              string
	DatabaseDSN          string
	AccrualSystemAddress string
}

func NewConfig() *Config {
	return &Config{}
}

func (o *Config) ParseConfig() {
	flag.StringVar(&o.RunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&o.DatabaseDSN, "d", "postgresql://localhost:54321/gofermart_db", "database address to connect")
	flag.StringVar(&o.AccrualSystemAddress, "r", "http://localhost:9090", "accrual system address to connect")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		o.RunAddr = envRunAddr
	}
	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		o.DatabaseDSN = envDatabaseDSN
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		o.AccrualSystemAddress = envAccrualSystemAddress
	}
}
