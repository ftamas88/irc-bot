package irc

import (
	"crypto/tls"
	"fmt"

	irc "github.com/fluffle/goirc/client"
)

func newClient(a args) *irc.Conn {
	cfg := irc.NewConfig(a.connect.Nick)

	if a.connect.SSL {
		cfg.SSL = true
		cfg.SSLConfig = &tls.Config{
			ServerName: a.connect.Server,
		}
	}

	cfg.Server = fmt.Sprintf("%s:%d", a.connect.Server, a.connect.Port)
	cfg.NewNick = func(n string) string {
		return n + "^"
	}
	client := irc.Client(cfg)

	return client
}
