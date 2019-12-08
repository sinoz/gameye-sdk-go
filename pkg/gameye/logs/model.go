package logs

// Entry represents an entry within the logs.
type Entry struct {
	LineKey int    `json:"lineKey"`
	Payload string `json:"payload"`
}
