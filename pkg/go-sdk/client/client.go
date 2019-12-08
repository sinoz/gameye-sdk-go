package client

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	heartbeat = 10 * 1000
)

// ErrMissingConfigField occurs when a config field is missing
var (
	ErrMissingConfigField = errors.New("missing field in config")
)

// GameyeClientConfig configures the GameyeClient
type GameyeClientConfig struct {
	Endpoint string
	Token    string
}

// GameyeClient is a simple wrapper for the gameye api, please use
// NewGameyeClient to create an instance.
type GameyeClient struct {
	config     GameyeClientConfig
	httpClient *http.Client
}

// NewGameyeClient constructs a new GameyeClient.
func NewGameyeClient(config GameyeClientConfig) (*GameyeClient, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	httpClient := &http.Client{}

	return &GameyeClient{
		config:     config,
		httpClient: httpClient,
	}, nil
}

func (config *GameyeClientConfig) validate() error {
	if strings.TrimSpace(config.Endpoint) == "" {
		config.Endpoint = os.Getenv("GAMEYE_API_ENDPOINT")
	}

	if strings.TrimSpace(config.Endpoint) == "" {
		return ErrMissingConfigField
	}

	if strings.TrimSpace(config.Token) == "" {
		config.Token = os.Getenv("GAMEYE_API_TOKEN")
	}

	if strings.TrimSpace(config.Token) == "" {
		return ErrMissingConfigField
	}

	return nil
}

func getCommandHeaders(config GameyeClientConfig) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.Token),
	}
}

func getStreamHeaders(config GameyeClientConfig) map[string]string {
	return map[string]string{
		"Authorization":        fmt.Sprintf("Bearer %s", config.Token),
		"Accept":               "application/x-ndjson",
		"x-heartbeat-interval": fmt.Sprintf("%d", heartbeat),
	}
}
