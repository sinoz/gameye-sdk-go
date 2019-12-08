package statistics

// State represents all the raw statistics data
type State struct {
	Statistics map[string]interface{}
}

// NewStateWithStatistics constructs a State object with a set of raw statistics
func NewStateWithStatistics(statistics map[string]interface{}) *State {
	return &State{Statistics: statistics}
}
