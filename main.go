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
	var actualCopy Bracket

	// TODO: goroutines here
	for _, player := range group.Players {
		player.Bracket = NewBracket(group.Field, player.Picks)
		if player.Name == "Actual" {
			actual = player.Bracket
			actualCopy = *actual.Copy()
		}
		if player.Name == "Linder" {
			fmt.Printf("\n\n\n %s's Bracket \n\n\n", player.Name)
			player.Bracket.PrettyPrint(os.Stdout, "\t\t\t")
		}
	}

	fmt.Printf("\n\n\n Actual Bracket \n\n\n")
	actual.PrettyPrint(os.Stdout, "\t\t\t")

	fmt.Printf("\n\n\n Actual filled with linder's picks Bracket \n\n\n")
	actualCopy.PrettyPrint(os.Stdout, "\t\t\t")

	fmt.Println(actualCopy.Points(actual, group.Rounds))
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

type Picks map[string]string

type Rounds map[int]Round

type Round struct {
	Name    string `json:"name"`
	Points  int    `json:"points"`
	AddSeed bool   `json:"add_seed"`
}
