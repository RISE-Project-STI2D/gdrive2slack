package gdrive2discord

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"../google"
	"../mailchimp"
)

type Configuration struct {
	BindAddress      string                     `json:"bindAddress"`
	Workers          int                        `json:"workers"`
	Interval         int                        `json:"interval"`
	GoogleTrackingId string                     `json:"googleTrackingId"`
	Google           *google.OauthConfiguration `json:"google"`
	Mailchimp        *mailchimp.Configuration   `json:"mailchimp"`
}

func LoadConfiguration(filename string) (*Configuration, error) {
	var self = new(Configuration)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(self)
	if err != nil {
		return nil, err
	}
	return self, nil
}

type Environment struct {
	Version         string
	Configuration   *Configuration
	Logger          *Logger
	HttpClient      *http.Client
	RegisterChannel chan *SubscriptionAndAccessToken
	SignalsChannel  chan os.Signal
}

func NewEnvironment(version string, conf *Configuration, logger *Logger) *Environment {
	e := &Environment{
		Version:       version,
		Configuration: conf,
		Logger:        logger,
		HttpClient: &http.Client{
			Timeout: time.Duration(15) * time.Second,
		},
		RegisterChannel: make(chan *SubscriptionAndAccessToken, 50),
		SignalsChannel:  make(chan os.Signal, 1),
	}
	signal.Notify(e.SignalsChannel, syscall.SIGINT, syscall.Signal(0xf))
	return e
}
