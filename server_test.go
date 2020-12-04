package main

import (
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
)


type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {

	store := StubPlayerStore{
		map[string]int{
			"Pepper" : 20,
			"Floyd" : 10,
		},
		nil,
	}
	server := &PlayerServer{&store}

	t.Run("Return scores for Pepper", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"
		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Return scores for Floyd", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"
		assertResponseBody(t, got, want)
		assertStatus(t, response.Code, http.StatusOK)
	})
	
	t.Run("Return 404 for a missing player", func(t *testing.T) {
		request := newGetScoreRequest("Lenny")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)

	})
}

func TestPOSTPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
    }
	server := &PlayerServer{&store}
	
	t.Run("it returns accepted POST request", func(t *testing.T) {
		request := newPostScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
	})

}

func TestScoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
	}
	
	server := &PlayerServer{&store}

	t.Run("it records win for a player when POST", func(t *testing.T) {

		player := "Pepper"
		request := newPostScoreRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

        if len(store.winCalls) != 1 {
            t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
		
		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}


func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newPostScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %q want %q", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("Got %d want %d", got ,want)
	}
}