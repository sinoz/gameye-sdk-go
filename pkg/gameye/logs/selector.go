package logs

import (
	"fmt"
	"sort"
	"strconv"
)

// SelectAll returns all logs from the store.
func SelectAll(state *State) []Entry {
	entries := []Entry{}
	rawLogs := state.Logs["line"]
	if rawLogs != nil {
		for _, v := range rawLogs.(map[string]interface{}) {
			rawLine := v.(map[string]interface{})
			if rawLine != nil {
				lineKeyInt, _ := strconv.Atoi(fmt.Sprintf("%v", rawLine["lineKey"]))
				logLine := Entry{
					LineKey: lineKeyInt,
					Payload: rawLine["payload"].(string),
				}
				entries = append(entries, logLine)
			}
		}
	}

	sort.Slice(entries, func(p, q int) bool {
		return (entries[q].LineKey - entries[p].LineKey) > 0
	})

	return entries
}

// SelectSince returns logs after the given line number.
func SelectSince(state *State, lineNumber int) []Entry {
	entries := []Entry{}
	rawLogs := state.Logs["line"]
	if rawLogs != nil {
		for _, v := range rawLogs.(map[string]interface{}) {
			rawLine := v.(map[string]interface{})
			if rawLine != nil {
				lineKeyInt, _ := strconv.Atoi(fmt.Sprintf("%v", rawLine["lineKey"]))
				logLine := Entry{
					LineKey: lineKeyInt,
					Payload: rawLine["payload"].(string),
				}
				if logLine.LineKey > lineNumber {
					entries = append(entries, logLine)
				}
			}
		}
	}

	sort.Slice(entries, func(p, q int) bool {
		return (entries[q].LineKey - entries[p].LineKey) > 0
	})

	return entries
}
