package app

import (
	"context"
	"github.com/ftamas88/irc-bot/internal/irc"
	"github.com/ftamas88/irc-bot/internal/tracker"
	"time"
)

// App struct defines parameters for the application.
type App struct {
	trackers *tracker.Service
	irc      *irc.Client
	timeout  time.Duration
}

// New function initiates and returns App struct.
func New(tr *tracker.Service, ic *irc.Client, dur time.Duration) App {
	return App{
		trackers: tr,
		irc:      ic,
		timeout:  dur,
	}
}

// Start handles the application startup and the graceful shutdown logic.
func (a App) Start(ctx context.Context) error {
	// shutdownChan channel is for shutting down all active connections gracefully.
	// Upon completion, the app gets notified and shuts itself down
	shutdownChan := make(chan struct{})

	go shutdownHandler(ctx, shutdownChan, a)

	// Run the IRC client
	go a.irc.Run(ctx)

	// Waiting for the shutdown to finish then inform the main goroutine.
	<-shutdownChan

	return nil
}

// shutdownHandler function runs as a goroutine behind the scene. It handles
// graceful application shutdown functionality. As soon as the context is
// cancelled with context.CancelFunc function, shutdown kicks in. It allows
// certain duration to active connections in order for them to finish their job
// first before terminating them. If any more requests come in after starting the
// shutdown process, they all will be refused. All the services that need to be
// shutdown should be listed in this function.
func shutdownHandler(ctx context.Context, shutdownChan chan<- struct{}, a App) {
	// Listening on parent context to see if any shutdown request was issued.
	// If so, following lines are executed otherwise this is a blocking line.
	<-ctx.Done()

	// Create a new context with a timeout duration. When it does time out, all
	// the associated resources are released.
	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	// Close shutdownChan channel to complete application shutdown.
	close(shutdownChan)
}
