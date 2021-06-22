package elite

type Rank int

const (
	None Rank = iota - 1
	Rookie
	Agent
	Officer
	SeniorOfficer
	Leader
)

var stringRanks = map[Rank]string{
	None:          "None",
	Rookie:        "Rookie",
	Agent:         "Agent",
	Officer:       "Officer",
	SeniorOfficer: "SeniorOfficer",
	Leader:        "Leader",
}

func (r Rank) ToString() string {
	return stringRanks[r]
}
