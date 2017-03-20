package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"sort"
)

var bracket_order = map[int]int {
	1: 0,
	16: 1,
	8: 2,
	9: 3,
	5: 4,
	12: 5,
	4: 6,
	13: 7,
	6: 8,
	11: 9,
	3: 10,
	14: 11,
	7: 12,
	10: 13,
	2: 14,
	15: 15,
}

func main () {

	fieldData, err := ioutil.ReadFile("./field.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var field FieldData
	json.Unmarshal(fieldData, &field)

	var regions = make(map[string]Bracket)

	for _, reg := range field.Field.Regions {
		sort.Sort(reg)
		var b = NewBracket(reg.Teams)
		regions[reg.Name] = b
	}

	for region, bracket := range regions {
		fmt.Printf("\n\n\n%s Bracket\n\n", region)
		bracket.PrettyPrint(os.Stdout, "\t\t\t")
	}
}

type Team struct {
	Name string `json:"name"`
	Region string `json:"region"`
	Seed int `json:"seed"`
	Eliminated string `json:"eliminated"`
}

type Region struct {
	Name string `json:"name"`
	Seed int `json:"seed"`
	Teams []Team `json:"teams"`
}

func (r Region) Len() int {
	return len(r.Teams)
}
func (r Region) Swap(i, j int) {
	r.Teams[i], r.Teams[j] = r.Teams[j], r.Teams[i]
}
func (r Region) Less(i, j int) bool {
	return bracket_order[r.Teams[i].Seed] < bracket_order[r.Teams[j].Seed]
}

type Field struct {
	Regions []Region`json:"regions"`
}

type FieldData struct {
	Field Field `json:"field"`
}