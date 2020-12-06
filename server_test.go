package main

import (
	"encoding/json"
	models "fantasy_league/Models"
	helpers "fantasy_league/TestHelpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []models.Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []models.Player {
	return s.league
}

func TestGETPlayers(t *testing.T) {

	store := StubPlayerStore{
		map[string]int{
			"Arthur": 20,
			"Dutch":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("Return scores for Arthur", func(t *testing.T) {
		request := helpers.NewGetScoreRequest("Arthur")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"
		helpers.AssertResponseBody(t, got, want)
		helpers.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Return scores for Dutch", func(t *testing.T) {
		request := helpers.NewGetScoreRequest("Dutch")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"
		helpers.AssertResponseBody(t, got, want)
		helpers.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Return 404 for a missing player", func(t *testing.T) {
		request := helpers.NewGetScoreRequest("Lenny")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertStatus(t, response.Code, http.StatusNotFound)

	})
}

func TestPOSTPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted POST request", func(t *testing.T) {
		request := helpers.NewPostScoreRequest("Micah")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		helpers.AssertStatus(t, response.Code, http.StatusAccepted)
	})

}

func TestScoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}

	server := NewPlayerServer(&store)

	t.Run("it records win for a player when POST", func(t *testing.T) {

		player := "Arthur"
		request := helpers.NewPostScoreRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		helpers.AssertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {

	wantedLeague := []models.Player{
		{"Arthur", 20},
		{"Dutch", 10},
		{"Lenny", 30},
	}

	store := StubPlayerStore{
		map[string]int{},
		nil,
		wantedLeague,
	}

	server := NewPlayerServer(&store)

	t.Run("It should return 200 for the /league endpoint", func(t *testing.T) {
		request := helpers.NewGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []models.Player

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}
		helpers.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("It should return the league table as a JSON", func(t *testing.T) {
		request := helpers.NewGetLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := helpers.ParseLeagueFromResponse(t, response.Body)
		helpers.AssertStatus(t, response.Code, http.StatusOK)
		helpers.AssertLeague(t, got, wantedLeague)
	})
}
