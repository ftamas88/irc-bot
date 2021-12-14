package tracker

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (t *Service) readTrackersConfig() error {
	var files []string

	if err := filepath.Walk(t.cfg.Folder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("error during reading the tracker files: %w", err)
	}

	for _, file := range files {
		xmlFile, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
		}

		var tr Tracker

		byteValue, _ := ioutil.ReadAll(xmlFile)
		err = xml.Unmarshal(byteValue, &tr)
		if err != nil {
			fmt.Printf("unable to parse the file: [%s]: %s", file, err.Error())
			continue
		}

		// TODO: optimise this, cache the loop
		for _, v := range t.config.Trackers {
			if v.Name == tr.Type {
				tr.Config = v
			}
		}

		t.trackers = append(t.trackers, tr)

		_ = xmlFile.Close()
	}

	if len(t.trackers) == 0 {
		return fmt.Errorf("error: %s", "unable to load any trackers")
	}

	return nil
}
