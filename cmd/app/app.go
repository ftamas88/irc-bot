package main

import (
	"context"
	"github.com/ftamas88/irc-bot/internal/app"
	"github.com/ftamas88/irc-bot/internal/config"
	"github.com/ftamas88/irc-bot/internal/irc"
	"github.com/ftamas88/irc-bot/internal/tracker"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// Log as Text
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		TimestampFormat:        "15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: false,
	})

	log.SetOutput(os.Stdout)

	// Available options in this app: Debug, Info, Warn
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("[~] (~‾▿‾)~ >> iRC downloader - initialising ")

	if err := godotenv.Load(); err != nil {
		log.Fatalf("[~] iRC downloader - Error loading .env file")
	}

	c, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("[~] iRC downloader - fatal error: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	// signChan channel is used to transmit signal notifications.
	signChan := make(chan os.Signal, 1)

	go signalHandler(signChan, cancel)

	trs, err := tracker.NewTrackerService(c.Tracker)
	if err != nil {
		log.Fatalf("[~] iRC downloader - fatal error: %s", err.Error())
	}

	ic := irc.NewClient(ctx, c, trs)

	if err := app.New(trs, ic, 3).Start(ctx); err != nil {
		log.Fatal(err)
	}
}

// signalHandler function runs as a goroutine behind the scene. It helps
// to trigger application shutdown by listening on certain signals as soon as
// they arrive. signChan channel is used to transmit signal notifications.
// Currently, handled signals are listed below as follows.
//
// os.Interrupt: Ctrl-C
// syscall.SIGTERM: kill PID, docker stop, docker down
func signalHandler(signChan chan os.Signal, cancel context.CancelFunc) {
	// Catch and relay certain signal(s) to signChan channel.
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	// Waiting for a signal to be sent on the signChan channel so that the
	// application shutdown can be triggered. If so, following lines are
	// executed otherwise this is a blocking line.
	sig := <-signChan

	log.Infof("[!] ¯`(>▂<)´¯ shutdown started with %v signal", sig)

	// Tell application to start abandoning its work.
	cancel()
}
