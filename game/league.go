package game

import (
	"io"
	"fmt"
	"encoding/json"
	models "fantasy_league/Models"
)

type League []models.Player

func NewLeague(rdr io.Reader) ([]models.Player, error) {
    var league []models.Player
    err := json.NewDecoder(rdr).Decode(&league)
    if err != nil {
        err = fmt.Errorf("problem parsing league, %v", err)
    }

    return league, err
}

func(l League) Find(name string) *models.Player {
    for i, p := range l {
        if p.Name == name {
            return &l[i]
        }
    }
    return nil
}
