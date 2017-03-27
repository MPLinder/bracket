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

	return construct([]Bracket{seed}, leaves[2:], picks)
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

	return construct([]Bracket{seed}, regions[2:], picks)
}

func (b *Bracket) Depth() int {
	if b.left == nil {
		return 1
	} else {
		return b.left.Depth() + 1
	}
}

// Round returns the integer value of the round in which a Bracket takes place
// assuming a 64 team field
func (b *Bracket) Round() int {
	return b.Depth() - 1
}

func (b *Bracket) Leaves() []Team {
	var result []Team
	if (b.left.value == Team{}) {
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

func gamePoints(bracket Bracket, actual Bracket, round Round) int {
	if (actual.value == Team{} || actual.value != bracket.value) {
		return 0
	}
	var adder = 0
	if round.AddSeed {
		adder = bracket.value.Seed
	}
	return round.Points + adder
}

func construct(parents []Bracket, leaves []Bracket, picks Picks) Bracket {

	if len(leaves) >= 2 {
		var newSeed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1], picks)}
		parents = append(parents, newSeed)

		var bracket = construct(parents, leaves[2:], picks)
		return bracket
	} else if len(leaves) == 2 {
		var newSeed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1], picks)}
		parents = append(parents, newSeed)

		return construct(parents, []Bracket{}, picks)
	} else {
		if len(parents) > 1 {
			var newSeed = Bracket{left: &parents[0], right: &parents[1], value: winner(parents[0], parents[1], picks)}
			return construct([]Bracket{newSeed}, parents[2:], picks)
		}
		return parents[0]
	}
}

func winner(left Bracket, right Bracket, picks Picks) Team {
	var winner = Team{}

	if picks[left.value.Name] == right.value.Name {
		winner = right.value
	} else if picks[right.value.Name] == left.value.Name {
		winner = left.value
	}

	return winner
}
