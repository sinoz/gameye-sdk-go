package session

// SelectSessionListForGame returns Sessions for a given gameKey
func SelectSessionListForGame(state *State, gameKey string) []Session {
	sessions := []Session{}
	for _, session := range state.Sessions {
		if session.Image == gameKey {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

// SelectSessionList returns all sessions currently active
func SelectSessionList(state *State) []Session {
	sessions := []Session{}
	for _, session := range state.Sessions {
		sessions = append(sessions, session)
	}
	return sessions
}

// SelectSession returns a session for the given sessionID
func SelectSession(state *State, sessionID string) Session {
	return state.Sessions[sessionID]
}
