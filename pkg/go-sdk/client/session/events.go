package session

import "encoding/json"

// Initialized is received when a session event stream is opened
type Initialized struct {
	Type    string                         `json:"type"`
	Payload SessionInitializedEventPayload `json:"payload"`
}

// SessionInitializedEventPayload is the payload for the SessionInitializedEvent
type SessionInitializedEventPayload struct {
	Sessions []SessionInitialized `json:"sessions"`
}

// SessionInitialized contains the information about an active sessions
type SessionInitialized struct {
	ID       string           `json:"id"`
	Image    string           `json:"image"`
	Location string           `json:"location"`
	Host     string           `json:"host"`
	Created  int64            `json:"created"`
	Port     map[string]int64 `json:"port"`
}

// Started is received when a session is started
type Started struct {
	Type    string                     `json:"type"`
	Payload SessionStartedEventPayload `json:"payload"`
}

// SessionStartedEventPayload is the payload for the SessionStartedEvent
type SessionStartedEventPayload struct {
	Session SessionStarted `json:"session"`
}

// SessionStarted contains the information about a started session
type SessionStarted struct {
	ID       string           `json:"id"`
	Image    string           `json:"image"`
	Location string           `json:"location"`
	Host     string           `json:"host"`
	Created  int64            `json:"created"`
	Port     map[string]int64 `json:"port"`
}

// Stopped is received when a session is ended
type Stopped struct {
	Type    string                     `json:"type"`
	Payload SessionStoppedEventPayload `json:"payload"`
}

// SessionStoppedEventPayload is the payload for the SessionStoppedEvent
type SessionStoppedEventPayload struct {
	Session SessionStopped `json:"session"`
}

// SessionStopped contains the identifier of a session
type SessionStopped struct {
	ID string `json:"id"`
}

// UnionEvent is used for partial deserialization
type UnionEvent struct {
	Type    string           `json:"type"`
	Payload *json.RawMessage `json:"payload"`
}
