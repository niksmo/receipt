package config

import "os"

type Config struct {
	LogLevel string
	Addr     string
	Login    string
	Password string
	SMTPHost string
	SMTPPort string
}

func Load() Config {
	return Config{
		LogLevel: os.Getenv("RECEIPT_LOG_LEVEL"),
		Addr:     os.Getenv("RECEIPT_ADDR"),
		Login:    os.Getenv("RECEIPT_LOGIN"),
		Password: os.Getenv("RECEIPT_PASSWORD"),
		SMTPHost: os.Getenv("RECEIPT_SMTP_HOST"),
		SMTPPort: os.Getenv("RECEIPT_SMTP_PORT"),
	}
}
