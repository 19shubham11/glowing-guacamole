package main

import (
	"log"
	"os"
	"net/http"
	game "fantasy_league/Game"
)

const dbFileName = "game.db.json"
const filePermissions = 0666

func main() {
	file, fileCreationError := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	
	if fileCreationError != nil {
		log.Fatalf("could not open file %s , %v", dbFileName, fileCreationError)
	}

	store, storeInitError := game.NewFileSystemPlayerStore(file)

	if storeInitError != nil {
		log.Fatalf("problem creating file system player store, %v ", storeInitError)
	}

	server := game.NewPlayerServer(store)

	if serverError := http.ListenAndServe(":5000", server); serverError != nil {
		log.Fatalf("could not listen on port 5000 %v", serverError)
	}
}
