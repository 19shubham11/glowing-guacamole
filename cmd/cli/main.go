package main

import (
	"log"
	"os"
	"io"
	"bufio"
	"fmt"
	"strings"
	game "fantasy_league/Game"
)
const dbFileName = "game.db.json"

type CLI struct{
	playerStore game.PlayerStore
	in *bufio.Scanner
}


func NewCLI(store game.PlayerStore, in io.Reader) *CLI {
    return &CLI{
        playerStore: store,
        in:          bufio.NewScanner(in),
    }
}

func(cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}


func(cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}


func main() {
	fmt.Println("Let's play poker!")
	fmt.Println("Type {Name} wins to record a win!")

	db, fileInitErr := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if fileInitErr != nil {
		log.Fatalf("Error opening file %v", fileInitErr)
	}

	
	store, storeInitError := game.NewFileSystemPlayerStore(db)

	if storeInitError !=nil {
		log.Fatalf("Error initialising player store %v", storeInitError)
	}

	cli := NewCLI(store, os.Stdin)
	cli.PlayPoker()
}
