package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"bytes"
)

func main() {
	groupData, err := ioutil.ReadFile("data/group.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var group Group
	err = json.Unmarshal(groupData, &group)
	if err != nil {
		fmt.Println("Unable to unmarshal group data: %v\n", err)
		os.Exit(1)
	}

	var actual Player
	actual.Bracket = NewBracket(group.Field, group.Actual.Picks)
	//actual.Bracket.PrettyPrint(os.Stdout, "\t\t\t")

	fmt.Printf("Round: %s\n", group.Rounds[actual.Bracket.LastCompleteRound()+1].Name)
	fmt.Printf("Number of teams left: %d\n", len(actual.Bracket.Leaves()))

	// TODO: goroutines here
	for i := range group.Players {
		group.Players[i].Bracket = NewBracket(group.Field, group.Players[i].Picks)
	}

	fmt.Println("Calculating all possible outcomes.")
	var allPossiblePicks []Picks
	switch actual.Bracket.LastCompleteRound() {
	case 0:
		fmt.Printf("It's way to early to do this. Chill for a bit.\n")
		return
	default:
		allPossiblePicks = actual.Bracket.AllPossiblePicks(ioutil.Discard, 2)
	}

	var allPossibleBrackets = []Bracket{}

	var possible Bracket
	for _, picks := range allPossiblePicks {
		possible = *actual.Bracket.Copy()
		possible.FillFromPicks(picks)
		allPossibleBrackets = append(allPossibleBrackets, possible)

		//// TODO: this only does one more round
		//var subPossiblePicks = possible.AllPossiblePicks(os.Stdout, 100)
		//var subPossible Bracket
		//for _, subPicks := range subPossiblePicks {
		//	subPossible = *possible.Copy()
		//	subPossible.FillFromPicks(subPicks)
		//	allPossibleBrackets = append(allPossibleBrackets, subPossible)
		//}
	}
	fmt.Println("\nDone.\n")

	fmt.Printf("Number of different ways you could have filled out a bracket for this round alone: %d\n\n", len(allPossibleBrackets))

	fmt.Println("Current Points:")

	sort.Slice(group.Players, func(i, j int) bool {
		return group.Players[i].Bracket.Points(actual.Bracket, group.Rounds) > group.Players[j].Bracket.Points(actual.Bracket, group.Rounds)
	})

	for i, player := range group.Players {
		fmt.Printf("\t %d. %s: %d\n", i+1, player.Name, player.Bracket.Points(actual.Bracket, group.Rounds))
	}

	var prefixBase = actual.Bracket.LastCompleteRound()
	var prefix string

	for _, bracket := range allPossibleBrackets {
		prefix = ""
		for i := 0; i < bracket.LastCompleteRound()-prefixBase; i++ {
			prefix += "\t"
		}

		//fmt.Printf("%sWhat will happen if these teams win: %v\n", prefix, bracket.RecentWinners())
		sort.Slice(group.Players, func(i, j int) bool {
			return group.Players[i].Bracket.Points(bracket, group.Rounds) > group.Players[j].Bracket.Points(bracket, group.Rounds)
		})

		var playerPossibleRound PlayerPossibleRound
		var points int
		for i, p := range group.Players {
			points = p.Bracket.Points(bracket, group.Rounds)
			//fmt.Printf("%s\t%s: %d\n", prefix, p.Name, points)
			playerPossibleRound = PlayerPossibleRound{Rank: i + 1, Points: points}
			p.PlayerPossibleRounds = append(p.PlayerPossibleRounds, playerPossibleRound)
		}

	}

	sort.Slice(group.Players, func(i, j int) bool {
		return group.Players[i].Name < group.Players[j].Name
	})

	for _, p := range group.Players {
		fmt.Println("\n", p.Name)
		fmt.Printf("\tPicks: %v\n", fmt.Sprint(strTeamSlice(p.Bracket.RoundWinners(actual.Bracket.LastCompleteRound()+1))))
		sort.Slice(p.PlayerPossibleRounds, func(i, j int) bool {
			if p.PlayerPossibleRounds[i].Rank == p.PlayerPossibleRounds[j].Rank {
				return p.PlayerPossibleRounds[i].Points > p.PlayerPossibleRounds[j].Points
			}
			return p.PlayerPossibleRounds[i].Rank < p.PlayerPossibleRounds[j].Rank
		})
		fmt.Printf("\tBest possible rank after this round: %d (%d points)\n", p.PlayerPossibleRounds[0].Rank, p.PlayerPossibleRounds[0].Points)
		fmt.Printf("\tWorst possible rank after this round: %d (%d points)\n", p.PlayerPossibleRounds[len(p.PlayerPossibleRounds)-1].Rank, p.PlayerPossibleRounds[len(p.PlayerPossibleRounds)-1].Points)
		fmt.Println("\n\tHow likely you are to be any particular rank:")

		var chance = make(map[int]int)
		for _, ppr := range p.PlayerPossibleRounds {
			chance[ppr.Rank] += 1
		}

		var ranks []int
		for c := range chance {
			ranks = append(ranks, c)
		}
		sort.Ints(ranks)

		//for rank, occurrences := range chance {
		for _, r := range ranks {
			fmt.Printf("\t\t%.2f%% chance of being ranked %d (%d of the %d scenarios)\n", (float64(chance[r]) / float64(len(p.PlayerPossibleRounds))) * 100, r, chance[r], len(p.PlayerPossibleRounds))
		}
	}
}

type Group struct {
	Field   Field     `json:"field"`
	Players []*Player `json:"players"`
	Actual  Player    `json:"actual"`
	Rounds  Rounds    `json:"rounds"`
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
	Name                 string `json:"name"`
	Picks                Picks  `json:"picks"`
	Bracket              Bracket
	PlayerPossibleRounds []PlayerPossibleRound
}

type PlayerPossibleRound struct {
	Rank   int
	Points int
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

func strTeamSlice(teams []Team) string {
	var b bytes.Buffer
	for i, t := range teams {
		if i != len(teams)-1 {
			fmt.Fprintf(&b, "%s, ", t)
		} else {
			fmt.Fprintf(&b, "%s", t)

		}
	}
	return b.String()
}
