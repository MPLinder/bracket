package main

import (
	"strconv"
	"io"
)

func NewBracket(teams []Team) Bracket {
	var leaves []Bracket

	for _, team := range teams {
		var leaf = Bracket{value: team}
		leaves = append(leaves, leaf)
	}

	var seed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1])}

	return construct([]Bracket{seed}, leaves[2:])
}

type Bracket struct {
	value Team
	left  *Bracket
	right *Bracket
}

func (b *Bracket) Depth() int {
	if (b.left == nil) {
		return 1
	} else {
		return b.left.Depth() + 1
	}
}

func (b *Bracket) String() string {
	var name string
	if b.value.Name != "" {
		var prefix string
		if (b.value.Eliminated != "") {
			prefix = "(X)"
		} else {
			prefix = ""
		}
		name = prefix + strconv.Itoa(b.value.Seed) + " " + b.value.Name
	} else {
		name = "________"
	}
	return name
}

func (t *Bracket) PrettyPrint(w io.Writer, prefix string) {
	var inner func(int, *Bracket)
	inner = func(depth int, child *Bracket) {
		if (child.left != nil) {
			inner(depth - 1, child.left)
		}
		for i := 0; i < depth; i++ {
			io.WriteString(w, prefix)
		}
		io.WriteString(w, child.String()+"\n")

		if (child.right != nil) {
			inner(depth - 1, child.right)
		}
	}
	inner(4, t)
}

func (b *Bracket) Leaves () []Team {
	var result []Team
	if (b.left.value == Team{}) {
		result = b.left.Leaves()
		result = append(result, b.right.Leaves()...)
	} else {
		return []Team{b.left.value, b.right.value}
	}

	return result
}

func construct(parents []Bracket, leaves []Bracket) Bracket {

	if (len(leaves) >= 2) {
		var newSeed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1])}
		parents = append(parents, newSeed)

		var bracket = construct(parents, leaves[2:])
		return bracket
	} else if (len(leaves) == 2) {
		var newSeed = Bracket{left: &leaves[0], right: &leaves[1], value: winner(leaves[0], leaves[1])}
		parents = append(parents, newSeed)

		return construct(parents, []Bracket{})
	} else {
		if (len(parents) > 1) {
			var newSeed = Bracket{left: &parents[0], right: &parents[1], value: winner(parents[0], parents[1])}
			return construct([]Bracket{newSeed}, parents[2:])
		}
		return parents[0]
	}
}

func winner(left Bracket, right Bracket) Team {
	var winner = Team{}

	if left.value.Eliminated == right.value.Name {
		winner = right.value
	} else if right.value.Eliminated == left.value.Name {
		winner = left.value
	}

	return winner
}