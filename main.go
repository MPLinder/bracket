package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
)

func main () {

	groupData, err := ioutil.ReadFile("./group.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	var group Group
	json.Unmarshal(groupData, &group)

	var bracket = NewBracket(group.Field)
	bracket.PrettyPrint(os.Stdout, "\t\t\t")
}

type Team struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	Seed         int `json:"seed"`
	EliminatedBy string `json:"eliminated_by"`
}

type Region struct {
	Name string `json:"name"`
	Seed int `json:"seed"`
	Teams []Team `json:"teams"`
}

type Field struct {
	Regions []Region`json:"regions"`
}

type Group struct {
	Field Field `json:"field"`
}