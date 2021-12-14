package tracker

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type TrackersConfig struct {
	Trackers []Config `yaml:"trackers"`
}

type Config struct {
	Name        string    `yaml:"name"`
	Enabled     bool      `yaml:"enabled" env-default:"false"`
	Debug       bool      `yaml:"debug"`
	Passkey     string    `yaml:"passkey"`
	AuthKey     string    `yaml:"auth_key"`
	Command     string    `yaml:"command"`
	DownloadDir string    `yaml:"download_dir"`
	Server      string    `yaml:"server"`
	SSL         bool      `yaml:"ssl"`
	Port        int       `yaml:"port"`
	Nick        string    `yaml:"nick"`
	MinFilesize string    `yaml:"min_filesize"`
	MaxFilesize string    `yaml:"max_filesize"`
	Filters     []Filters `yaml:"filters,omitempty"`
}

type Filters struct {
	Categories Categories `env:"categories"`
}

type Categories map[string]struct{}

func (c *Categories) UnmarshalYAML(value *yaml.Node) error {
	var tmpCategories string
	if err := value.Decode(&tmpCategories); err != nil {
		return err
	}

	cats := map[string]struct{}{}
	for _, value := range strings.Split(tmpCategories, ",") {
		cats[value] = struct{}{}
	}

	*c = cats

	return nil
}
