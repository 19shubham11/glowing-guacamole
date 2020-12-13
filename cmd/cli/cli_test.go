package main

import (
	"strings"
	game "github.com/19shubham11/glowing-guacamole/game"
	"testing"
)


type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   game.League
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() game.League {
	return s.league
}


func assertPlayerWins(t *testing.T, stubPlayerStore *StubPlayerStore, got, want string) {
	t.Helper()

	if len(stubPlayerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}


func TestCLI(t *testing.T) {
	t.Run("record win for dutch", func(t *testing.T){
		in := strings.NewReader("dutch wins\n")
		playerStore := &StubPlayerStore{}

		cli := NewCLI(playerStore, in)
		cli.PlayPoker()

		got := playerStore.winCalls[0]
		want := "dutch"

		assertPlayerWins(t, playerStore, got, want)
	})

	t.Run("record win for arthur", func(t *testing.T){
		in := strings.NewReader("arthur wins\n")
		playerStore := &StubPlayerStore{}

		cli := NewCLI(playerStore, in)
		cli.PlayPoker()

		got := playerStore.winCalls[0]
		want := "arthur"

		assertPlayerWins(t, playerStore, got, want)
	})
}
