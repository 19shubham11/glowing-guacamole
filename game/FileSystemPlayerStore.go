package game

import (
	"sort"
	"os"
	"encoding/json"
	models "github.com/19shubham11/glowing-guacamole/models"
	"fmt"
)


type tape struct{
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}

type FileSystemPlayerStore struct {
	database *json.Encoder
	league League
}


func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)
	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting stats of the file %v", err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0,0)
	}
	return err
}



func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0,0)

	fileInitErr := initialisePlayerDBFile(file)
	
	if fileInitErr != nil {
		return nil, fmt.Errorf("problems initialising empty file %v", fileInitErr)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem parsing store from file %v", err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league: league,
	}, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool{
		return f.league[i].Wins > f.league[j].Wins
	})
    return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins ++
	}

	if player == nil {
		f.league = append(f.league, models.Player{name, 1})
	}
	f.database.Encode(f.league)
}
