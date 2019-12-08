package gameye

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Gameye/gameye-sdk-go/pkg/gameye/logs"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/match"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/patch"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/session"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/statistics"
	"github.com/Gameye/messaging-client-go/pkg/command"
	"github.com/Gameye/messaging-client-go/pkg/eventstream"
)

const (
	heartbeat = 10 * 1000
)

// ErrMissingConfigField occurs when a config field is missing
var (
	ErrMissingConfigField = errors.New("missing field in config")
)

// ClientConfig configures the gameye Client.
// TODO use functional options pattern to support default values in a clean way:
// - https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type ClientConfig struct {
	Endpoint string
	Token    string
}

// Client is a simple wrapper for the gameye api, please use
// NewClient to create an instance.
type Client struct {
	config     ClientConfig
	httpClient *http.Client
}

// NewClient constructs a new gameye Client.
func NewClient(config ClientConfig) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	httpClient := &http.Client{}

	return &Client{
		config:     config,
		httpClient: httpClient,
	}, nil
}

// StartMatch starts a new match with the given paramters
func (client *Client) StartMatch(
	matchKey string,
	gameKey string,
	locationKeys []string,
	templateKey string,
	config map[string]interface{},
	endCallbackURL string,
) error {
	action := match.Start{
		Type: "start-match",
		Payload: match.StartPayload{
			MatchKey:       matchKey,
			GameKey:        gameKey,
			LocationKeys:   locationKeys,
			TemplateKey:    templateKey,
			Config:         config,
			EndCallbackURL: endCallbackURL,
		},
	}

	url := fmt.Sprintf("%s/command/%s", client.config.Endpoint, action.Type)

	return command.Invoke(
		url,
		action.Payload,
		getCommandHeaders(client.config),
	)
}

// StopMatch ends a match with the given matchKey
func (client *Client) StopMatch(matchKey string) error {
	action := match.Stop{
		Type: "stop-match",
		Payload: match.StopPayload{
			MatchKey: matchKey,
		},
	}

	url := fmt.Sprintf("%s/command/%s", client.config.Endpoint, action.Type)

	return command.Invoke(
		url,
		action.Payload,
		getCommandHeaders(client.config),
	)
}

// SubscribeStatisticsEvents adds a subscriber to statistic events
func (client *Client) SubscribeStatisticsEvents(matchKey string, onStateChanged func(*statistics.State)) (err error) {
	ctx := context.Background()
	url := fmt.Sprintf("%s/query/statistic", client.config.Endpoint)

	statistics.SubscribeState("client.statistics.internal", onStateChanged)
	decoder, err := eventstream.Create(
		ctx,
		url,
		map[string]string{"matchKey": matchKey},
		getStreamHeaders(client.config),
	)

	if err != nil {
		return err
	}

	go func() {
		for {
			var patches []patch.Patch
			err = decoder.Decode(&patches)

			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
				break
			} else if patches != nil && len(patches) > 0 {
				statistics.Dispatch(patches)
			}
		}

		statistics.UnsubscribeState("client.statistics.internal")
	}()

	return nil
}

// SubscribeSessionEvents adds a subscriber to session events
func (client *Client) SubscribeSessionEvents(onStateChanged func(*session.State)) error {
	ctx := context.Background()
	url := fmt.Sprintf("%s/query/session", client.config.Endpoint)

	session.SubscribeState("client.session.internal", onStateChanged)

	decoder, err := eventstream.Create(ctx, url, nil, getStreamHeaders(client.config))
	if err != nil {
		return err
	}

	go func() {
		for {
			var action session.UnionEvent
			err := decoder.Decode(&action)

			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
				break
			} else {
				session.Dispatch(&action)
			}
		}

		session.UnsubscribeState("client.session.internal")
	}()

	return nil
}

// SubscribeLogEvents adds a subscriber to log events
func (client *Client) SubscribeLogEvents(matchKey string, onStateChanged func(*logs.State)) (err error) {
	ctx := context.Background()
	url := fmt.Sprintf("%s/query/log", client.config.Endpoint)

	logs.SubscribeState("client.log.internal", onStateChanged)
	decoder, err := eventstream.Create(
		ctx,
		url,
		map[string]string{"matchKey": matchKey},
		getStreamHeaders(client.config),
	)

	if err != nil {
		return err
	}

	go func() {
		for {
			var patches []patch.Patch
			err = decoder.Decode(&patches)

			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
				break
			} else if patches != nil && len(patches) > 0 {
				logs.Dispatch(patches)
			}
		}

		logs.UnsubscribeState("client.log.internal")
	}()

	return nil
}

func (config *ClientConfig) validate() error {
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

func getCommandHeaders(config ClientConfig) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.Token),
	}
}

func getStreamHeaders(config ClientConfig) map[string]string {
	return map[string]string{
		"Authorization":        fmt.Sprintf("Bearer %s", config.Token),
		"Accept":               "application/x-ndjson",
		"x-heartbeat-interval": fmt.Sprintf("%d", heartbeat),
	}
}
