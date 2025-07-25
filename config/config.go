package config

type Config struct {
	LogLevel string
	Addr     string
	Login    string
	Password string
	SMTPHost string
	SMTPPort string
}

func Load() Config {
	return Config{}
}
