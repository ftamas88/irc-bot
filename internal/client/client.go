package client

import (
	"crypto/tls"
	"fmt"

	irc "github.com/fluffle/goirc/client"
	"github.com/ftamas88/irc-bot/internal/config"
)

func Client(config *config.Config) *irc.Conn {
	cfg := irc.NewConfig(config.Nick)

	if config.Port != 6667 {
		cfg.SSL = true
		cfg.SSLConfig = &tls.Config{ServerName: config.Server}
	}

	cfg.Server = fmt.Sprintf("%s:%d", config.Server, config.Port)
	cfg.NewNick = func(n string) string {
		return n + "^"
	}
	client := irc.Client(cfg)

	return client
}
