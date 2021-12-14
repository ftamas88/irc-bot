package irc

import (
	"fmt"
	"github.com/dustin/go-humanize"
	irc "github.com/fluffle/goirc/client"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Torrent contains the basic information retrieved from the channel
type Torrent struct {
	BaseURL     string
	DownloadURL string
	Category    string
	Name        string
	Size        string
	ID          int
}

type Size struct {
	Size float64
	Unit string
}

func handleMessage(a args, line *irc.Line) {
	// Remove the unicode/non printable characters which screws up the regexp..
	message := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, line.Text())

	re := regexp.MustCompile(a.tracker.ParseInfo.LinePatterns.Extract.Regex.Value)
	match := re.FindStringSubmatch(message)

	/*
		<linepatterns>
					<extract>
						<regex value="\[NEW TORRENT in .\d{0,}?(\D{1,}).*]\d{0,}\s?(.*)\14?\s>\d{1,}? {0,}?\d?.* (\d.*.*\DiB).*in.*>\s{1,}(https?:\/\/ncore.pro.*action=).*id=(\d+)\s?"/>
						<vars>
							<var name="category"/>
							<var name="torrentName"/>
							<var name="$torrentSize"/>
							<var name="$baseUrl"/>
							<var name="$torrentId"/>
						</vars>
					</extract>
				</linepatterns>
				<linematched>
					<var name="torrentUrl">
						<var name="$baseUrl"/>
						<string value="download&amp;id="/>
						<var name="$torrentId"/>
						<string value="&amp;key="/>
						<var name="passkey"/>
					</var>
				</linematched>
	*/

	var t Torrent
	if len(match) > 0 {
		for k, value := range a.tracker.ParseInfo.LinePatterns.Extract.Vars.Var {
			switch value.Name {
			case "category":
				t.Category = match[(k + 1)]
				break
			case "torrentName":
				t.Name = match[(k + 1)]
			case "$torrentSize":
				t.Size = match[(k + 1)]
			case "$baseUrl":
				t.BaseURL = match[(k + 1)]
			case "$torrentId":
				t.ID, _ = strconv.Atoi(match[(k + 1)])
			default:
				break
			}
		}

		fs, err := humanize.ParseBytes(t.Size)
		if err != nil {
			log.Warnf("unable to parse filesize: %s", err.Error())
			return
		}

		minFs, _ := humanize.ParseBytes(a.tracker.Config.MinFilesize)
		maxFs, _ := humanize.ParseBytes(a.tracker.Config.MaxFilesize)

		inFilter := false
		for _, v := range a.tracker.Config.Filters {
			_, categoryFound := v.Categories[t.Category]
			if len(v.Categories) > 0 && !categoryFound {
				log.Debugf("[-] Skipping torrent, category out of range: %s", t.Category)
				return
			}
			if fs >= minFs && fs <= maxFs {
				inFilter = true
			}
		}

		if !inFilter {
			log.Debugf("[-] Skipping torrent, filesize/filters doesn't match: %s", t.Size)

			return
		}

		log.
			Debugf("[*] New torrent [*]\nName:\t\t%s\nID:\t\t%d\nSize:\t\t%s\nCategory:\t%s", t.Name, t.ID, t.Size, t.Category)

		if err := downloadFile(a, t); err != nil {
			log.Warnf("[~] iRC downloader - ERROR: unable to download file: %s E: %s", t.Name, err.Error())
		}

		return
	}

	// Unknown message
	log.
		WithField("message", message).
		WithField("date", line.Time.Format("15:04:05")).
		WithField("nick", line.Nick).
		Info("Unknown message")

	return
}

// DownloadFile will download an url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(a args, t Torrent) error {
	fPath := fmt.Sprintf("%s/%d.torrent", a.tracker.Config.DownloadDir, t.ID)

	// Create the download link
	downloadLink := ""
	for _, v := range a.tracker.ParseInfo.LineMatched.Var {
		if v.Name == "torrentUrl" {
			for i, j := range v.Var {
				downloadLink += j.Name
				if len(v.String) > i {
					downloadLink += v.String[i].Value
				}
			}
		}
	}

	downloadLink = strings.Replace(downloadLink, "$baseUrl", t.BaseURL, -1)
	downloadLink = strings.Replace(downloadLink, "$torrentId", strconv.Itoa(t.ID), -1)
	downloadLink = strings.Replace(downloadLink, "passkey", a.tracker.Config.Passkey, -1)
	downloadLink = strings.Replace(downloadLink, "auth_key", a.tracker.Config.AuthKey, -1)

	// Get the data
	resp, err := http.Get(downloadLink)
	if err != nil {
		log.Warnf("error downloading the file: %s", err.Error())
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Create the file
	out, err := os.Create(fPath)
	if err != nil {
		log.Warnf("error downloading the file: %s", err.Error())
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Warnf("error downloading the file: %s", err.Error())
		return err
	}

	log.
		WithField("Category", t.Category).
		WithField("Size", t.Size).
		Infof("[~] iRC downloader - %s", t.Name)

	return nil
}
