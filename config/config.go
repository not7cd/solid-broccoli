package config

import (
	"github.com/kelseyhightower/envconfig"
)

const PrefixCmd = '.'

type Query byte

const (
	QueryAdd     Query = '+'
	QueryRandom  Query = '.' // no arg
	QueryInspect Query = '?' // no arg
)

type Config struct {
	DiscordToken string `required:"true"  split_words:"true"`
	SqlitePath   string `required:"true"  split_words:"true"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}

	// TODO: Load from TOML

	if err := envconfig.Process("BROCCOLI", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
