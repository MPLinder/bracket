package main


//type Region struct {
//	Name string
//	Seed int
//	Teams []Team
//}
//
//type Game struct {
//	Team1 Team
//	Team2 Team
//	Winner Team
//}
//
//type Player struct {
//	Name string
//	Rank int
//	Points int
//}
//
//type Round struct {
//	Name string
//}
//
//type Bracket struct {
//	Regions []Region
//}


//type Team struct {
//	Name string
//	Seed int
//	Eliminated bool
//}
//
//type Bracket struct {
//	bp BracketProvider
//}
//
//type Region struct {
//	name string
//	seed int
//	tp TeamProvider
//}
//
//type BracketProvider interface {
//	Regions()
//}
//
//type TeamProvider interface {
//	Teams(reg Region) []Team
//	TeamsAilve(reg Region) []Team
//}

//type Team struct {
//	Name string `json:"name"`
//	Region string `json:"region"`
//	Seed int `json:"seed"`
//	Eliminated bool `json:"eliminated"`
//}
//
//type Region struct {
//	Name string `json:"name"`
//	Seed int `json:"seed"`
//}
//
//type Field struct {
//	Regions []Region`json:"regions"`
//	Teams   []Team `json:"teams"`
//}
//
//type FieldData struct {
//	Field Field `json:"field"`
//}
//
//type TeamProvider interface {
//	Teams(Field) []Team
//}