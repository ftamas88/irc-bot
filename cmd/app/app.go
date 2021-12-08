package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	irc "github.com/fluffle/goirc/client"
	"github.com/ftamas88/irc-bot/internal/client"
	"github.com/ftamas88/irc-bot/internal/config"
	log "github.com/sirupsen/logrus"
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

// Torrent contains the basic information retrieved from the channel
type Torrent struct {
	Category string
	Name     string
	Size     Size
	ID       int
}

type Size struct {
	Size float64
	Unit string
}

func main() {
	log.Info("[~] iRC downloader - initalising")

	config, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("[~] iRC downloader - fatal error: %s", err.Error())
	}

	client := client.Client(config)

	// Initial stuff
	client.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			log.
				WithField("server", config.Server).
				WithField("port", config.Port).
				Info("[+] Connected to the iRC server")

			conn.Privmsg("NBOT", fmt.Sprintf("!invite %s", config.InviteCode))
			conn.Mode("#ncore-bot", "+r")
			conn.Join("#ncore-bot")
		},
	)

	// On Receive
	client.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {
			handleMessage(config, line)
		},
	)

	// Disconnect
	quit := make(chan bool)
	client.HandleFunc(
		irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			quit <- true
		},
	)

	// Connect
	if err := client.Connect(); err != nil {
		log.Infof("Connection error: %s\n", err.Error())
	}

	go shutdownHandler(client)

	// Wait for disconnect
	<-quit
}

func handleMessage(cfg *config.Config, line *irc.Line) {
	/*
		https://regex101.com/
	*/
	nCoreRegexp := `\[NEW TORRENT in .\d{0,}?(\D{1,}).*]\d{0,}\s?(.*)\14?\s>\d{1,}? {0,}?(\d{1,5}\.?\d{0,2}) (MiB|GiB|TiB).*in.*>\s{1,}https:\/\/[a-zA-Z{2,}].*id=(\d+)\s?`

	// Remove the unicode/non printable characters which screws up the regexp..
	message := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, line.Text())

	re := regexp.MustCompile(nCoreRegexp)
	match := re.FindStringSubmatch(message)

	if len(match) > 0 {
		id, _ := strconv.Atoi(match[5])
		size, _ := strconv.ParseFloat(match[3], 64)
		t := Torrent{
			ID:       id,
			Name:     match[2],
			Category: match[1],
			Size: Size{
				Size: size,
				Unit: match[4],
			},
		}

		// Create the download link
		downloadLink := strings.Replace(cfg.DownloadLink, "[ID]", strconv.Itoa(t.ID), -1)
		downloadLink = strings.Replace(downloadLink, "[PASSKEY]", cfg.Passkey, -1)

		log.
			Debugf("[*] New torrent [*]\nName:\t\t%s\nID:\t\t%d\nSize:\t\t%.02f %s\nCategory:\t%s", t.Name, t.ID, t.Size.Size, t.Size.Unit, t.Category)

		// Download the .torrent file
		go func(dir string, torrent *Torrent) {
			if err := downloadFile(
				fmt.Sprintf("%s/%s.torrent", dir, torrent.Name),
				downloadLink,
			); err != nil {
				log.Warnf("[~] iRC downloader - ERROR: unable to download file: %s E: %s", torrent.Name, err.Error())
			}
			log.
				WithField("Category", torrent.Category).
				WithField("Size", fmt.Sprintf("%.2f %s", torrent.Size.Size, torrent.Size.Unit)).
				Infof("[~] iRC downloader - %s", torrent.Name)
		}(cfg.DownloadDir, &t)

		return
	}

	// Unknown message
	log.Infof("%s >> %s >> %v\n", line.Time.Format("15:04:05"), line.Nick, message)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Warnf("error downloading the file: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Warnf("error downloading the file: %s", err.Error())
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Warnf("error downloading the file: %s", err.Error())
		return err
	}

	return nil
}

// shutdownHandler listens for a SIGTERM signal
// and gracefully cancels the main application context
// once this is completed exits the app
func shutdownHandler(client *irc.Conn) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Warnf("[!] iRC downloader - Bye :o/")
	client.Quit("bye")

	<-time.After(time.Duration(1) * time.Second)
	os.Exit(1)
}
