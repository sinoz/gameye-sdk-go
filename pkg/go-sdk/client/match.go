package client

import (
	"fmt"

	"github.com/Gameye/gameye-sdk-go/pkg/go-sdk/client/match"
	"github.com/Gameye/messaging-client-go/pkg/command"
)

// StartMatch starts a new match with the given paramters
func StartMatch(client *GameyeClient,
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
func StopMatch(client GameyeClient, matchKey string) error {
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
