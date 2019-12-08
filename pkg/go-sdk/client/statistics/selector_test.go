package statistics

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Gameye/gameye-sdk-go/pkg/go-sdk/helper"
)

func readStats() (result map[string]interface{}) {
	result, err := helper.ReadFileAsJSON("./testdata/stats.json")
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func TestSelectsAllPlayers(test *testing.T) {
	state := NewStateWithStatistics(readStats())
	filtered := SelectPlayerList(state)

	assert.Equal(test, 10, len(filtered))
	for _, v := range filtered {
		if v.PlayerKey == "11" {
			assert.Equal(test, "Ivan", v.Name)
			return
		}
	}
}

func TestSelectsOnlyRequestedPlayers(test *testing.T) {
	state := NewStateWithStatistics(readStats())
	filtered, err := SelectPlayerListForTeam(state, "1")

	assert.Nil(test, err)
	assert.Equal(test, 5, len(filtered))
	assertions := 0
	for _, v := range filtered {
		if v.PlayerKey == "7" {
			assertions++
		}

		if v.Name == "Seth" {
			assertions++
		}

		if v.Name == "Zane" {
			assertions++
			test.Fail()
		}
		if v.PlayerKey == "8" {
			assertions++
			test.Fail()
		}
	}

	assert.Equal(test, 2, assertions)

	found, err := SelectPlayer(state, "9")
	assert.Nil(test, err)
	assert.Equal(test, "Vladimir", found.Name)
}

func TestSelectsAllTeams(test *testing.T) {
	state := NewStateWithStatistics(readStats())
	filtered := SelectTeamList(state)

	assert.Equal(test, 2, len(filtered))
	for _, v := range filtered {
		if v.TeamKey == "2" {
			assert.Equal(test, "Terrorists", v.Name)
			return
		}
	}
}
func TestSelectsOnlyRequestedTeams(test *testing.T) {
	state := NewStateWithStatistics(readStats())
	team, err := SelectTeam(state, "1")

	assert.Nil(test, err)
	assert.NotEmpty(test, team)
	assert.Equal(test, "1", team.TeamKey)
	assert.Equal(test, "Counter Terrorists", team.Name)
}

func TestSelectsRounds(test *testing.T) {
	state := NewStateWithStatistics(readStats())
	rounds := SelectRounds(state)

	assert.NotEmpty(test, rounds)
	assert.Equal(test, 2, rounds)
}

func TestSelectsRawStatistics(test *testing.T) {
	state := NewStateWithStatistics(readStats())
	raw, _ := SelectRawStatistics(state)

	unmarshaled := make(map[string]interface{})
	data := readStats()

	err := json.Unmarshal([]byte(raw), &unmarshaled)

	assert.Nil(test, err)
	assert.Equal(test, data, unmarshaled)
}
