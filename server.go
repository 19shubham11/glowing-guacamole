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
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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


func GetPlayerScore(name string) int {
	if name == "Pepper" {
		return 20
	}

	if name == "Floyd" {
		return 10
	}
	return 0
}

func getPlayerName(url string) string {
	return strings.TrimPrefix(url, "/players/")
}