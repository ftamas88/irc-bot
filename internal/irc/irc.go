package irc

import (
	"context"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/ftamas88/irc-bot/internal/config"
	"github.com/ftamas88/irc-bot/internal/tracker"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	ctx context.Context
	cfg *config.Config
	trs *tracker.Service
}

type args struct {
	connect config.IRC
	tracker tracker.Tracker
}

func NewClient(ctx context.Context, cfg *config.Config, t *tracker.Service) *Client {
	return &Client{
		ctx: ctx,
		cfg: cfg,
		trs: t,
	}
}

func (c *Client) Run(ctx context.Context) {
	for _, v := range c.trs.Trackers() {
		if !v.Config.Enabled {
			continue
		}
		log.
			WithField("server", v.LongName).
			WithField("Host", v.Config.Server).
			WithField("Port", v.Config.Port).
			WithField("Nick", v.Config.Nick).
			WithField("Channel", v.Servers.Server[0].Channels).
			Info("Connecting to IRC server")

		param := args{
			connect: config.IRC{
				Server: v.Config.Server,
				Port:   v.Config.Port,
				Nick:   v.Config.Nick,
				SSL:    v.Config.SSL,
			},
			tracker: v,
		}

		// Connect to the server
		go func(a args) {
			go c.handle(a)
		}(param)
	}

	// Waiting for context cancel
	<-ctx.Done()
}

func (c *Client) handle(a args) {
	cl := newClient(a)

	// Initial stuff
	cl.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			log.
				WithField("server", a.connect.Server).
				WithField("port", a.connect.Port).
				Info("[+] Connected to the iRC server")

			ch := fmt.Sprintf("%s", a.tracker.Servers.Server[0].Channels)

			if a.tracker.Config.Command != "" {
				conn.Raw(a.tracker.Config.Command)
			}

			conn.Mode(ch, "+r")
			conn.Join(ch)
		},
	)

	// On Receive
	cl.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			handleMessage(a, line)
		},
	)

	// Disconnect
	cl.HandleFunc(
		irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			panic("disconnected from IRC ¯\\_( ͡° ͜ʖ ͡°)_/¯") // TODO: handle reconnect logic
		},
	)

	// Connect
	if err := cl.Connect(); err != nil {
		log.Infof("Connection error: %s\n", err.Error())
	}
}
