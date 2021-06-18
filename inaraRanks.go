package main

type InaraRank int

const (
	Outsider InaraRank = iota - 1
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

type EliteRank int

const (
	None EliteRank = iota - 1
	Rookie
	Agent
	Officer
	SeniorOfficer
	Leader
)

var RankMap map[InaraRank]EliteRank = map[InaraRank]EliteRank{
	Outsider:                None,
	Recruit:                 Rookie,
	Reserve:                 Agent,
	CoPilot:                 Agent,
	Pilot:                   Agent,
	Wingman:                 Agent,
	SeniorWingman:           Agent,
	Veteran:                 Agent,
	FlightLeader:            Officer,
	FlightGroupLeader:       Officer,
	OperationsOfficer:       SeniorOfficer,
	ChiefOfStaff:            SeniorOfficer,
	DeputySquadronCommander: SeniorOfficer,
	SquadronCommander:       Leader,
}

func (i InaraRank) GetEliteRank() EliteRank {
	return RankMap[i]
}
