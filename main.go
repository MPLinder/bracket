package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

func main() {
	groupData, err := ioutil.ReadFile("./group.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	var group Group
	json.Unmarshal(groupData, &group)

	var actual Player
	actual.Bracket = NewBracket(group.Field, group.Actual.Picks)
	//actual.Bracket.PrettyPrint(os.Stdout, "\t\t\t")

	// TODO: goroutines here
	for i := range group.Players {
		group.Players[i].Bracket = NewBracket(group.Field, group.Players[i].Picks)
	}

	sort.Slice(group.Players, func(i, j int) bool {
		return group.Players[i].Bracket.Points(actual.Bracket, group.Rounds) > group.Players[j].Bracket.Points(actual.Bracket, group.Rounds)
	})

	fmt.Println("Current Points:")
	for _, player := range group.Players {
		fmt.Printf("\t %s: %d\n", player.Name, player.Bracket.Points(actual.Bracket, group.Rounds))
	}

	var allPossiblePicks = actual.Bracket.AllPossiblePicks()
	var allPossibleBrackets = []Bracket{}

	var possible Bracket
	for _, picks := range allPossiblePicks {
		possible = *actual.Bracket.Copy()
		possible.FillFromPicks(picks)
		allPossibleBrackets = append(allPossibleBrackets, possible)

		// TODO: this only does one more round
		var subPossiblePicks = possible.AllPossiblePicks()
		var subPossible Bracket
		for _, subPicks := range subPossiblePicks {
			subPossible = *possible.Copy()
			subPossible.FillFromPicks(subPicks)
			allPossibleBrackets = append(allPossibleBrackets, subPossible)
		}
	}

	var prefixBase = allPossibleBrackets[0].LastCompleteRound()
	var prefix string
	for _, bracket := range allPossibleBrackets {
		prefix = ""
		for i := 0; i < bracket.LastCompleteRound()-prefixBase; i++ {
			prefix += "\t"
		}

		fmt.Printf("%sWinners: %v\n", prefix, bracket.RecentWinners())
		sort.Slice(group.Players, func(i, j int) bool {
			return group.Players[i].Bracket.Points(bracket, group.Rounds) > group.Players[j].Bracket.Points(bracket, group.Rounds)
		})

		for _, p := range group.Players {
			fmt.Printf("%s\t%s: %d\n", prefix, p.Name, p.Bracket.Points(bracket, group.Rounds))
		}
	}
}

type Group struct {
	Field   Field    `json:"field"`
	Players []Player `json:"players"`
	Actual  Player   `json:"actual"`
	Rounds  Rounds   `json:"rounds"`
}

type Team struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Seed   int    `json:"seed"`
}

func (t Team) String() string {
	return t.Name
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
