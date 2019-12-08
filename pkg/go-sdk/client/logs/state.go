package logs

// State represents all the raw log data
type State struct {
	Logs map[string]interface{}
}

// NewStateWithLogs constructs a State object with a set of raw logs
func NewStateWithLogs(logs map[string]interface{}) *State {
	return &State{Logs: logs}
}
