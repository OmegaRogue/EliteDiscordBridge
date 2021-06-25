//go:generate go-enum -f=$GOFILE --marshal --noprefix --sqlnullint -t ../assets/zerolog.gotmpl

package inara

import elite "github.com/OmegaRogue/eliteJournal"

// Rank is an enumeration of Inara Squadron Ranks
/*
ENUM(
outsider=-1
recruit
reserve
co-pilot
pilot
wingman
senior wingman
veteran
flight leader
flight group leader
operations officer
chief of staff
deputy squadron commander
squadron commander
)
*/
type Rank int

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

func (x Rank) GetEliteRank() elite.Rank {
	return rankMap[x]
}
