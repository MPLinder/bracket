package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	groupData, err := ioutil.ReadFile("./group.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var group Group
	json.Unmarshal(groupData, &group)

	var actual Bracket
	var linder Bracket

	// TODO: goroutines here
	for _, player := range group.Players {
		player.Bracket = NewBracket(group.Field, player.Picks)
		if player.Name == "Actual" {
			actual = player.Bracket
		}
		if player.Name == "Linder" {
			linder = player.Bracket
		}
	}

	var allPossiblePicks = actual.AllPossiblePicks()

	var possible Bracket
	for _, picks := range allPossiblePicks {
		possible = *actual.Copy()
		possible.FillFromPicks(picks)
		fmt.Println(linder.Points(possible, group.Rounds))
	}
}

type Group struct {
	Field   Field    `json:"field"`
	Players []Player `json:"players"`
	Rounds  Rounds   `json:"rounds"`
}

type Team struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Seed   int    `json:"seed"`
}

type Region struct {
	Name  string `json:"name"`
	Seed  int    `json:"seed"`
	Teams []Team `json:"teams"`
}

type Field struct {
	Regions []Region `json:"regions"`
}

type Player struct {
	Name    string `json:"name"`
	Picks   Picks  `json:"picks"`
	Bracket Bracket
}

// Picks is a mapping of team name to the integer round in which you expect them to be eliminated.
type Picks map[string]int

func (p *Picks) Copy() Picks {
	var newPicks = make(Picks)
	for k, v := range *p {
		newPicks[k] = v
	}
	return newPicks
}

type Rounds map[int]Round

type Round struct {
	Name    string `json:"name"`
	Points  int    `json:"points"`
	AddSeed bool   `json:"add_seed"`
}
