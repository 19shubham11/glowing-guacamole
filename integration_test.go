package main

import (
	models "fantasy_league/Models"
	helpers "fantasy_league/TestHelpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetreivingThem(t *testing.T) {

	database, cleandb := helpers.CreateTempFile(t, "")
	defer cleandb()

	store := &FileSystemPlayerStore{database, nil}
	server := NewPlayerServer(store)
	playerName := "lenny"
	server.ServeHTTP(httptest.NewRecorder(), helpers.NewPostScoreRequest(playerName))
	server.ServeHTTP(httptest.NewRecorder(), helpers.NewPostScoreRequest(playerName))
	server.ServeHTTP(httptest.NewRecorder(), helpers.NewPostScoreRequest(playerName))

	t.Run("Get Score", func(t *testing.T) {
		request := helpers.NewGetScoreRequest(playerName)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		helpers.AssertStatus(t, response.Code, http.StatusOK)
		helpers.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("Get League", func(t *testing.T) {
		request := helpers.NewGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		helpers.AssertStatus(t, response.Code, http.StatusOK)

		got := helpers.ParseLeagueFromResponse(t, response.Body)
		want := []models.Player{
			{Name: playerName, Wins: 3},
		}
		helpers.AssertLeague(t, got, want)
	})
}
