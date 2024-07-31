package tools

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Host string `env:"ONI_HOST, default=localhost"`
	Port string `env:"ONI_PORT, default=8080"`

	DB struct {
		Postgres struct {
			Database string `env:"ONI_DB_PG_DATABASE, default=oni"`
			Host     string `env:"ONI_DB_PG_HOST, default=postgres"`
			User     string `env:"ONI_DB_PG_USER, default=oni"`
			Password string `env:"ONI_DB_PG_PASSWORD, default=admin"`
		}
		SQLite struct {
			Path string `env:"ONI_DB_LITE_PATH, default=./oni.db"`
		}
	}

	Manga struct {
		Path string `env:"ONI_MANGA_PATH, default=./manga"`
	}
}

func LoadConfig() *Config {
	ctx := context.Background()
	var config Config

	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}
