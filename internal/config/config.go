package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       string
	Port         int
	Nick         string
	Passkey      string
	InviteCode   string
	DownloadDir  string
	DownloadLink string
}

func ReadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[~] iRC downloader - Error loading .env file")

		return nil, err
	}

	var err error
	var server, nick string
	var port int

	if os.Getenv("LOCAL") == "true" {
		server = os.Getenv("LOCAL_SERVER")
		port, err = strconv.Atoi(os.Getenv("LOCAL_PORT"))
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", err.Error())
		}
		nick = os.Getenv("LOCAL_NICK")
	} else {
		server = os.Getenv("SERVER")
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", err.Error())
		}
		nick = os.Getenv("NICK")
	}

	return &Config{
		Server:       server,
		Port:         port,
		Nick:         nick,
		Passkey:      os.Getenv("PASSKEY"),
		InviteCode:   os.Getenv("INVITE_CODE"),
		DownloadDir:  os.Getenv("DOWNLOAD_DIR"),
		DownloadLink: os.Getenv("DOWNLOAD_LINK"),
	}, nil
}
