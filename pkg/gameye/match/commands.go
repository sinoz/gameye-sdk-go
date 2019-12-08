package match

import "encoding/json"

// Start is a command to start a new match.
type Start struct {
	Type    string       `json:"type"`
	Payload StartPayload `json:"payload"`
}

// StartPayload is the carried payload for the Start command.
type StartPayload struct {
	AccountKey     string                 `json:"accountKey"`
	MatchKey       string                 `json:"matchKey"`
	GameKey        string                 `json:"gameKey"`
	LocationKeys   []string               `json:"locationKeys"`
	TemplateKey    string                 `json:"templateKey"`
	Config         map[string]interface{} `json:"config"`
	EndCallbackURL string                 `json:"endCallbackUrl"`
}

// Stop is a command to stop a currently ongoing match.
type Stop struct {
	Type    string      `json:"type"`
	Payload StopPayload `json:"payload"`
}

// StopPayload is the payload for the StopMatchCommand
type StopPayload struct {
	AccountKey string `json:"accountKey"`
	MatchKey   string `json:"matchKey"`
}

// UnionCommand is used for partial deserialization
type UnionCommand struct {
	Type    string           `json:"type"`
	Payload *json.RawMessage `json:"payload"`
}
