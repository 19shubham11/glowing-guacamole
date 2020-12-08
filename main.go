package main

import (
	"log"
	"os"
	"net/http"
)

const dbFileName = "game.db.json"
const filePermissions = 0666

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	
	if err != nil {
		log.Fatalf("could not open file %s , %v", dbFileName, err)
	}

	store := &FileSystemPlayerStore{db, nil}
	server := NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
