package main

import (
	"edDiscord/elite"
	"edDiscord/inara"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	InaraUser InaraUser
	Servers   []*Server  `gorm:"many2many:user_guilds;"`
	Roles     []*Role    `gorm:"many2many:user_roles;"`
	InaraRank inara.Rank `gorm:"default:null"`
	EliteRank elite.Rank `gorm:"default:null"`
}

type Server struct {
	gorm.Model
	Users []*User `gorm:"many2many:user_guilds;"`
	Roles []Role
}

type Role struct {
	gorm.Model
	Color     int
	Name      string
	Position  int
	Users     []*User `gorm:"many2many:user_roles;"`
	ServerID  uint
	InaraRank inara.Rank `gorm:"default:null"`
	EliteRank elite.Rank `gorm:"default:null"`
}

type InaraUser struct {
	gorm.Model
	UserID          uint
	Name            string
	CommanderName   string
	Allegiance      elite.Allegiance
	Power           elite.Power
	AvatarURL       string
	URL             string
	InaraWingID     uint
	InaraSquadronID uint

	CombatRank           int
	CombatProgress       float64
	TradeRank            int
	TradeProgress        float64
	ExplorationRank      int
	ExplorationProgress  float64
	CqcRank              int
	CqcProgress          float64
	SoldierRank          int
	SoldierProgress      float64
	ExobiologistRank     int
	ExobiologistProgress float64
	EmpireRank           int
	EmpireProgress       float64
	FederationRank       int
	FederationProgress   float64
}

type InaraWing struct {
	gorm.Model
	Name            string
	URL             string
	InaraSquadronID uint
	Users           []InaraUser
}

type InaraSquadron struct {
	gorm.Model
	Name  string
	URL   string
	Wings []InaraWing
	Users []InaraUser
}
