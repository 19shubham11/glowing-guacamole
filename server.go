package main

import (
	"fmt"
	"net/http"
	"strings"
)


type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

func NewPlayerServer(store PlayerStore) * PlayerServer {

	p := new(PlayerServer)
	p.store = store

	router := http.NewServeMux()

	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func(p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, r)
	case http.MethodGet:
		p.showScore(w, r)	
	}
}

func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := getPlayerName(r.URL.Path)
	w.WriteHeader(http.StatusAccepted)
	p.store.RecordWin(player)
	return
}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
	player := getPlayerName(r.URL.Path)
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}

func getPlayerName(url string) string {
	return strings.TrimPrefix(url, "/players/")
}
