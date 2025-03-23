package utils

import (
	"fmt"
	"net/url"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ListenAddress         string        `env:"LISTEN_ADDRESS" envDefault:":3000"`
	Prefork               bool          `env:"PREFORK" envDefault:"true"`
	DatabaseFilename      string        `env:"DATABASE_FILENAME" envDefault:"local.db"`
	DatabaseTursoUrl      url.URL       `env:"DATABASE_TURSO_URL"`
	DatabaseTursoToken    string        `env:"DATABASE_TURSO_TOKEN"`
	DatabaseTursoSyncTime time.Duration `env:"DATABASE_TURSO_SYNC_TIME" envDefault:"5m"`

	SmtpUsername string `env:"SMTP_USERNAME,required,notEmpty"`
	SmtpHost     string `env:"SMTP_HOST,required,notEmpty"`
	SmtpPassword string `env:"SMTP_PASSWORD,required,notEmpty"`
	SmtpAddress  string `env:"SMTP_ADDRESS,required,notEmpty"`
	SmtpPort     int    `env:"SMTP_PORT,required"`
}

func GetConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("Failed To Get Parse Config: %v\n", err)
	}
	return cfg
}
