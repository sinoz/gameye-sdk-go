package session

// State represents all the session data
type State struct {
	Sessions map[string]Session
}

// NewStateWithSessions constructs a State object with a set of sessions
func NewStateWithSessions(sessions map[string]Session) *State {
	return &State{sessions}
}
