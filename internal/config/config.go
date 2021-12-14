package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IRC
	Tracker
	Passkey      string  `env:"PASSKEY" env-required:"true"`
	InviteCode   string  `env:"INVITE_CODE" env-required:"true"`
	DownloadDir  string  `env:"DOWNLOAD_DIR" env-required:"true"`
	DownloadLink string  `env:"DOWNLOAD_LINK" env-required:"true"`
	MinFilesize  float64 `env:"MIN_FILESIZE" env-required:"true"`
	MaxFilesize  float64 `env:"MAX_FILESIZE" env-required:"true"`
}

type IRC struct {
	Server string `env:"SERVER"`
	Port   int    `env:"PORT"`
	Nick   string `env:"NICK"`
	SSL    bool

	Local bool `env:"LOCAL"`

	// Local version
	LocalServer string `env:"LOCAL_SERVER"`
	LocalPort   int    `env:"LOCAL_PORT"`
	LocalNick   string `env:"LOCAL_NICK"`
}

type Tracker struct {
	Folder   string `env:"TRACKERS_FOLDER" env-required:"true"`
	Enabled  []string
	Disabled []string
}

func ReadConfig() (*Config, error) {
	var err error
	var cfg Config

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error during setting up the config: %w", err)
	}

	if cfg.IRC.Local {
		cfg.IRC.Server = cfg.IRC.LocalServer
		cfg.IRC.Port = cfg.IRC.LocalPort
		cfg.IRC.Nick = cfg.IRC.LocalNick
	}

	return &cfg, nil
}
