package main

import (
	"io"
	"sort"
	"strconv"
)

var bracket_order = map[int]int{
	1:  0,
	16: 1,
	8:  2,
	9:  3,
	5:  4,
	12: 5,
	4:  6,
	13: 7,
	6:  8,
	11: 9,
	3:  10,
	14: 11,
	7:  12,
	10: 13,
	2:  14,
	15: 15,
}

type Bracket struct {
	value Team
	left  *Bracket
	right *Bracket
}

func NewRegion(region Region, picks Picks) Bracket {
	sort.Slice(region.Teams, func(i, j int) bool { return bracket_order[region.Teams[i].Seed] < bracket_order[region.Teams[j].Seed] })

	var teams = region.Teams
	var leaves []Bracket

	for _, team := range teams {
		var leaf = Bracket{value: team}
		leaves = append(leaves, leaf)
	}

	var seed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1], picks)}

	return constructFromSlice([]Bracket{seed}, leaves[2:], picks)
}

func NewBracket(field Field, picks Picks) Bracket {
	sort.Slice(field.Regions, func(i, j int) bool {
		return bracket_order[field.Regions[i].Seed] < bracket_order[field.Regions[j].Seed]
	})

	var regions []Bracket

	for _, reg := range field.Regions {
		var b = NewRegion(reg, picks)
		regions = append(regions, b)
	}

	var seed = Bracket{left: &regions[0], right: &regions[1], value: winner(regions[0], regions[1], picks)}

	return constructFromSlice([]Bracket{seed}, regions[2:], picks)
}

func (b Bracket) Copy() *Bracket {
	var ret = Bracket{value: b.value}
	if b.left == nil {
		return &ret
	}
	ret.left = b.left.Copy()
	ret.right = b.right.Copy()
	return &ret
}

func (b Bracket) Depth() int {
	if b.left == nil {
		return 1
	} else {
		return b.left.Depth() + 1
	}
}

func (b *Bracket) FillFromPicks(picks Picks) {
	if b.left == nil {
		return
	}

	if b.left.value == (Team{}) {
		b.left.FillFromPicks(picks)
	}

	if b.right.value == (Team{}) {
		b.right.FillFromPicks(picks)
	}

	b.value = winner(*b.left, *b.right, picks)
}

func (b *Bracket) AllPossiblePicks() []Picks {
	if b.value != (Team{}) {
		return []Picks{}
	}

	var ret []Picks
	var picks = make(Picks)

	ret = append(ret, picks)

	allPossiblePicksHelper(*b, &ret)
	return ret
}

// Round returns the integer value of the round in which a Bracket takes place
// assuming a 64 team field
func (b *Bracket) Round() int {
	return b.Depth() - 1
}

// LastCompleteRound returns the integer value of the last complete round in which a Bracket takes place
// assuming a 64 team field
func (b *Bracket) LastCompleteRound() int {
	if b.value != (Team{}) {
		return b.Round()
	}

	var left = b.left.LastCompleteRound()
	var right = b.right.LastCompleteRound()

	if left < right {
		return left
	}
	return right
}

func (b *Bracket) Leaves() []Team {
	var result []Team
	if b.left.value == (Team{}) {
		result = b.left.Leaves()
		result = append(result, b.right.Leaves()...)
	} else {
		return []Team{b.left.value, b.right.value}
	}

	return result
}

func (b *Bracket) String() string {
	var name string
	if b.value.Name != "" {
		name = strconv.Itoa(b.value.Seed) + " " + b.value.Name
	} else {
		name = "________"
	}
	return name
}

func (b *Bracket) PrettyPrint(w io.Writer, prefix string) {
	var inner func(int, *Bracket)
	inner = func(depth int, child *Bracket) {
		if child.left != nil {
			inner(depth-1, child.left)
		}
		for i := 0; i < depth; i++ {
			io.WriteString(w, prefix)
		}
		io.WriteString(w, child.String()+"\n")

		if child.right != nil {
			inner(depth-1, child.right)
		}
	}
	inner(b.Depth()-1, b)
}

func (b *Bracket) Points(actual Bracket, rounds Rounds) int {
	if b.Round() == 0 {
		return 0
	}
	return gamePoints(*b, actual, rounds[b.Round()]) + b.left.Points(*actual.left, rounds) + b.right.Points(*actual.right, rounds)
}

func (b *Bracket) RecentWinners() []Team {
	var recentWinners []Team
	recentWinnersHelper(b, &recentWinners)
	return recentWinners
}

func (b *Bracket) RoundWinners(round int) []Team {
	var roundWinners []Team
	roundWinnersHelper(b, round, &roundWinners)
	return roundWinners
}

func recentWinnersHelper(b *Bracket, rw *[]Team) {
	if b.value == (Team{}) {
		recentWinnersHelper(b.left, rw)
		recentWinnersHelper(b.right, rw)
	} else {
		*rw = append(*rw, b.value)
	}
	return
}

func roundWinnersHelper(b *Bracket, round int, rw *[]Team) {
	if b.Round() > round {
		roundWinnersHelper(b.left, round, rw)
		roundWinnersHelper(b.right, round, rw)
	} else {
		*rw = append(*rw, b.value)
	}
	return
}

func gamePoints(bracket Bracket, actual Bracket, round Round) int {
	if actual.value == (Team{}) || actual.value != bracket.value {
		return 0
	}
	var adder = 0
	if round.AddSeed {
		adder = bracket.value.Seed
	}
	return round.Points + adder
}

func constructFromSlice(parents []Bracket, leaves []Bracket, picks Picks) Bracket {

	if len(leaves) >= 2 {
		var newSeed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1], picks)}
		parents = append(parents, newSeed)

		var bracket = constructFromSlice(parents, leaves[2:], picks)
		return bracket
	} else if len(leaves) == 2 {
		var newSeed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1], picks)}
		parents = append(parents, newSeed)

		return constructFromSlice(parents, []Bracket{}, picks)
	} else {
		if len(parents) > 1 {
			var newSeed = Bracket{left: &parents[0], right: &parents[1], value: winner(parents[0], parents[1], picks)}
			return constructFromSlice([]Bracket{newSeed}, parents[2:], picks)
		}
		return parents[0]
	}
}

func winner(left Bracket, right Bracket, picks Picks) Team {
	var winner = Team{}

	if picks[left.value.Name] > picks[right.value.Name] {
		winner = left.value
	} else if picks[right.value.Name] > picks[left.value.Name] {
		winner = right.value
	}

	return winner
}

func allPossiblePicksHelper(bracket Bracket, picksSlice *[]Picks) {

	if bracket.value != (Team{}) {
		return
	}

	if bracket.left.value != (Team{}) && bracket.right.value != (Team{}) {
		var newPicks Picks
		var newPicksSlice []Picks

		// For each Picks in the slice of Picks, copy it and one team winning to the original and a different team winning to the copy
		// then append them back together
		for _, p := range *picksSlice {
			newPicks = p.Copy()

			p[bracket.left.value.Name] = bracket.Round()
			p[bracket.right.value.Name] = bracket.Round() + 1

			newPicks[bracket.left.value.Name] = bracket.Round() + 1
			newPicks[bracket.right.value.Name] = bracket.Round()

			newPicksSlice = append(newPicksSlice, newPicks)
		}

		*picksSlice = append(*picksSlice, newPicksSlice...)
	}

	allPossiblePicksHelper(*bracket.left, picksSlice)
	allPossiblePicksHelper(*bracket.right, picksSlice)

	return
}
