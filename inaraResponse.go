package main

type InaraStatus int

type InaraPilotRank struct {
	RankName     string  `json:"rankName"`
	RankValue    string  `json:"rankValue"`
	RankProgress float64 `json:"rankProgress"`
}

type InaraShip struct {
	ShipType  string `json:"shipType"`
	ShipName  string `json:"shipName"`
	ShipIdent string `json:"shipIdent"`
	ShipRole  string `json:"shipRole"`
}

type InaraSquadron struct {
	SquadronID           int    `json:"squadronID"`
	SquadronName         string `json:"squadronName"`
	SquadronMembersCount int    `json:"squadronMembersCount"`
	SquadronMemberRank   string `json:"squadronMemberRank"`
	InaraURL             string `json:"inaraURL"`
}

type InaraWing struct {
	WingID           int    `json:"wingID"`
	WingName         string `json:"wingName"`
	WingMembersCount int    `json:"wingMembersCount"`
	WingMemberRank   string `json:"wingMemberRank"`
	InaraURL         string `json:"inaraURL"`
}
type InaraProfile struct {
	UserID                  int              `json:"userID"`
	UserName                string           `json:"userName"`
	CommanderName           string           `json:"commanderName"`
	CommanderRanksPilot     []InaraPilotRank `json:"commanderRanksPilot"`
	PreferredAllegianceName string           `json:"preferredAllegianceName"`
	PreferredPowerName      string           `json:"preferredPowerName"`
	CommanderMainShip       InaraShip        `json:"commanderMainShip"`
	CommanderSquadron       InaraSquadron    `json:"commanderSquadron"`
	CommanderWing           InaraWing        `json:"commanderWing"`
	AvatarImageURL          string           `json:"avatarImageURL"`
	InaraURL                string           `json:"inaraURL"`
}

type InaraResponseData interface {
}

type InaraHeaderEventData struct {
	UserID   int    `json:"userID"`
	UserName string `json:"userName"`
}

type InaraResponseHeader struct {
	EventData   InaraHeaderEventData `json:"eventData"`
	EventStatus InaraStatus          `json:"eventStatus"`
}

type InaraResponse struct {
	Header InaraResponseHeader  `json:"header"`
	Events []InaraResponseEvent `json:"events"`
}

type InaraResponseEvent struct {
	EventCustomID int               `json:"eventCustomID"`
	EventStatus   int               `json:"eventStatus"`
	EventData     InaraResponseData `json:"eventData"`
}
