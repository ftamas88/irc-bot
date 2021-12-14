package tracker

import (
	"fmt"
	"github.com/ftamas88/irc-bot/internal/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type Service struct {
	cfg      config.Tracker
	trackers []Tracker
	config   TrackersConfig
}

func NewTrackerService(cfg config.Tracker) (*Service, error) {
	s := &Service{
		cfg:      cfg,
		trackers: []Tracker{},
		config:   TrackersConfig{},
	}

	err := cleanenv.ReadConfig("config/config.yaml", &s.config)
	if err != nil {
		return nil, fmt.Errorf("error during setting up the trackers config: %w", err)
	}

	if err := s.readTrackersConfig(); err != nil {
		return nil, fmt.Errorf("failed to parse the tracker config: %w", err)
	}

	return s, nil
}

func (t *Service) Trackers() []Tracker {
	return t.trackers
}
