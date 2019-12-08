package main

import (
	"fmt"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/logs"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/session"
	"github.com/Gameye/gameye-sdk-go/pkg/gameye/statistics"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
)

func main() {
	config := gameye.ClientConfig{
		Endpoint: "",
		Token:    "",
	}

	client, err := gameye.NewClient(config)
	handleErr(err)

	sessionID := uuid.New().String()

	var activeSession = make(chan session.Session)
	log.Printf("Subscribing to session events %v\n", sessionID)
	onSessionState := func(state *session.State) {
		foundSession := session.SelectSession(state, sessionID)
		if foundSession.ID != "" {
			log.Printf("Match Ready! %v", foundSession)
		}
		activeSession <- foundSession
	}

	err = client.SubscribeSessionEvents(onSessionState)
	handleErr(err)

	log.Printf("Starting Match %v\n", sessionID)
	err = client.StartMatch(
		sessionID,
		"csgo-dem",
		[]string{"frankfurt"},
		"bots",
		map[string]interface{}{"maxRounds": 2},
		"",
	)
	handleErr(err)

	currentLine := 0
	var allLogs []logs.Entry
	onLogState := func(state *logs.State) {
		newLogs := logs.SelectSince(state, currentLine)
		for _, v := range newLogs {
			log.Printf("%d: %s", v.LineKey, v.Payload)
		}
		currentLine += len(newLogs)
		allLogs = logs.SelectAll(state)
	}

	err = client.SubscribeLogEvents(sessionID, onLogState)
	handleErr(err)

	rawStats := ""
	onStatisticsState := func(state *statistics.State) {
		statistics.SelectPlayerList(state)
		rawStats, err = statistics.SelectRawStatistics(state)
		handleErr(err)
	}

	err = client.SubscribeStatisticsEvents(sessionID, onStatisticsState)
	handleErr(err)

	// Wait for a session to not be empty
	for currentSession := range activeSession {
		if currentSession.ID != "" {
			break
		}
	}

	log.Printf("Session found %v\n", sessionID)

	// Wait for a session to become empty
	for currentSession := range activeSession {
		if currentSession.ID == "" {
			break
		}
	}

	log.Printf("Match finished %v\n", sessionID)

	file, err := os.Create("logs.txt")
	if file != nil && allLogs != nil {
		for _, v := range allLogs {
			_, err = io.WriteString(file, fmt.Sprintf("%d: %s", v.LineKey, v.Payload))
		}

		err = file.Close()
		handleErr(err)
	}

	file, err = os.Create("stats.txt")
	if file != nil {
		_, err = io.WriteString(file, rawStats)
		err = file.Close()
		handleErr(err)
	}

}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
