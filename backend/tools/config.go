package tools

import (
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Host string
	Port string

	DB struct {
		Path string
	}

	Manga struct {
		Path string
	}
}

func LoadConfig() *Config {
	k := koanf.New(".")

	// Load config from environment variables
	k.Load(env.Provider("ONI_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "ONI_")), "_", ".", -1)
	}), nil)

	// Unmarshal the config into a struct
	var config *Config
	k.Unmarshal("", config)

	return config
}
