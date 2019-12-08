package logs

import (
	"log"
	"testing"

	"github.com/Gameye/gameye-sdk-go/pkg/go-sdk/helper"
	"github.com/stretchr/testify/assert"
)

func readLogs() (json map[string]interface{}) {
	json, err := helper.ReadFileAsJSON("./testdata/logs.json")
	if err != nil {
		log.Fatal(err)
	}
	return json
}

func TestSelectsAllLogs(test *testing.T) {
	state := NewStateWithLogs(readLogs())
	filtered := SelectAll(state)

	assert.Equal(test, 1095, len(filtered))
	assert.Equal(test, 561, filtered[560].LineKey)
	assert.Equal(test, "$L 09/16/2019 - 12:39:05: \"Joe<4><BOT><TERRORIST>\" dropped \"vesthelm\"", filtered[896].Payload)
}

func TestSelectsOnlyRequestedLogs(test *testing.T) {
	state := NewStateWithLogs(readLogs())
	filtered := SelectSince(state, 912)

	assert.Equal(test, 1095-912, len(filtered))
}
