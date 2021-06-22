package inara

type Status int

type SquadronID int
type WingID int
type UserID int

type PilotRank struct {
	RankName     string  `json:"rankName" `
	RankValue    string  `json:"rankValue" `
	RankProgress float64 `json:"rankProgress" `
}

type Ship struct {
	ShipType  string `json:"shipType" `
	ShipName  string `json:"shipName" `
	ShipIdent string `json:"shipIdent" `
	ShipRole  string `json:"shipRole" `
}

type Squadron struct {
	SquadronID           SquadronID `json:"squadronID" `
	SquadronName         string     `json:"squadronName" `
	SquadronMembersCount int        `json:"squadronMembersCount" `
	SquadronMemberRank   string     `json:"squadronMemberRank" `
	InaraURL             string     `json:"inaraURL" `
}

type Wing struct {
	WingID           WingID `json:"wingID"`
	WingName         string `json:"wingName" `
	WingMembersCount int    `json:"wingMembersCount" `
	WingMemberRank   string `json:"wingMemberRank" `
	InaraURL         string `json:"inaraURL" `
}
type Profile struct {
	UserID                  UserID      `json:"userID" `
	UserName                string      `json:"userName" `
	CommanderName           string      `json:"commanderName" `
	CommanderRanksPilot     []PilotRank `json:"commanderRanksPilot" `
	PreferredAllegianceName string      `json:"preferredAllegianceName" `
	PreferredPowerName      string      `json:"preferredPowerName" `
	CommanderMainShip       Ship        `json:"commanderMainShip" `
	CommanderSquadron       Squadron    `json:"commanderSquadron" `
	CommanderWing           Wing        `json:"commanderWing" `
	AvatarImageURL          string      `json:"avatarImageURL" `
	InaraURL                string      `json:"inaraURL" `
}

type ResponseData interface {
}

type HeaderEventData struct {
	UserID   UserID `json:"userID"`
	UserName string `json:"userName"`
}

type ResponseHeader struct {
	EventData   HeaderEventData `json:"eventData"`
	EventStatus Status          `json:"eventStatus"`
}

type Response struct {
	Header ResponseHeader  `json:"header"`
	Events []ResponseEvent `json:"events"`
}

type ResponseEvent struct {
	EventCustomID int     `json:"eventCustomID"`
	EventStatus   int     `json:"eventStatus"`
	EventData     Profile `json:"eventData"`
}
