package inara

import "edDiscord/elite"

type Rank int

const (
	Outsider Rank = iota - 1
	Recruit
	Reserve
	CoPilot
	Pilot
	Wingman
	SeniorWingman
	Veteran
	FlightLeader
	FlightGroupLeader
	OperationsOfficer
	ChiefOfStaff
	DeputySquadronCommander
	SquadronCommander
)

var rankStrings = map[string]Rank{
	"recruit":                   Recruit,
	"reserve":                   Reserve,
	"co-pilot":                  CoPilot,
	"pilot":                     Pilot,
	"wingman":                   Wingman,
	"senior wingman":            SeniorWingman,
	"veteran":                   Veteran,
	"flight leader":             FlightLeader,
	"flight group leader":       FlightGroupLeader,
	"operations officer":        OperationsOfficer,
	"chief of staff":            ChiefOfStaff,
	"deputy squadron commander": DeputySquadronCommander,
	"squadron commander":        SquadronCommander,
}
var stringRanks = map[Rank]string{
	Outsider:                "none",
	Recruit:                 "recruit",
	Reserve:                 "reserve",
	CoPilot:                 "co-pilot",
	Pilot:                   "pilot",
	Wingman:                 "wingman",
	SeniorWingman:           "senior wingman",
	Veteran:                 "veteran",
	FlightLeader:            "flight leader",
	FlightGroupLeader:       "flight group leader",
	OperationsOfficer:       "operations officer",
	ChiefOfStaff:            "chief of staff",
	DeputySquadronCommander: "deputy squadron commander",
	SquadronCommander:       "squadron commander",
}

func (r Rank) ToString() string {
	return stringRanks[r]
}

func GetInaraRank(rank string) Rank {
	return rankStrings[rank]
}

var rankMap = map[Rank]elite.Rank{
	Outsider:                elite.None,
	Recruit:                 elite.Rookie,
	Reserve:                 elite.Agent,
	CoPilot:                 elite.Agent,
	Pilot:                   elite.Agent,
	Wingman:                 elite.Agent,
	SeniorWingman:           elite.Agent,
	Veteran:                 elite.Agent,
	FlightLeader:            elite.Officer,
	FlightGroupLeader:       elite.Officer,
	OperationsOfficer:       elite.SeniorOfficer,
	ChiefOfStaff:            elite.SeniorOfficer,
	DeputySquadronCommander: elite.SeniorOfficer,
	SquadronCommander:       elite.Leader,
}

func (r Rank) GetEliteRank() elite.Rank {
	return rankMap[r]
}
