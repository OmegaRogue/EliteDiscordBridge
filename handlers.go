package main

import (
	"context"
	"strings"

	"edDiscord/inara"
	elite "github.com/OmegaRogue/eliteJournal"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func roleCreated(s *discordgo.Session, e *discordgo.GuildRoleCreate) {
	guildFlake, err := snowflake.ParseString(e.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	roleFlake, err := snowflake.ParseString(e.Role.ID)
	if err != nil {
		log.Fatalf("parse role snowflake: %+v", err)
	}
	db.FirstOrCreate(
		&Role{
			Model:    gorm.Model{ID: uint(roleFlake)},
			Color:    e.Role.Color,
			Name:     e.Role.Name,
			Position: e.Role.Position,
			ServerID: uint(guildFlake),
		},
	)
}

func roleUpdated(s *discordgo.Session, e *discordgo.GuildRoleUpdate) {
	guildFlake, err := snowflake.ParseString(e.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	roleFlake, err := snowflake.ParseString(e.Role.ID)
	if err != nil {
		log.Fatalf("parse role snowflake: %+v", err)
	}
	db.Updates(
		&Role{
			Model:    gorm.Model{ID: uint(roleFlake)},
			Color:    e.Role.Color,
			Name:     e.Role.Name,
			Position: e.Role.Position,
			ServerID: uint(guildFlake),
		},
	)

}

func guildCreated(s *discordgo.Session, e *discordgo.GuildCreate) {
	guildFlake, err := snowflake.ParseString(e.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	db.FirstOrCreate(&Server{Model: gorm.Model{ID: uint(guildFlake)}})
	for _, role := range e.Roles {
		roleFlake, err := snowflake.ParseString(role.ID)
		if err != nil {
			log.Fatalf("parse role snowflake: %+v", err)
		}
		db.FirstOrCreate(
			&Role{
				Model:    gorm.Model{ID: uint(roleFlake)},
				Color:    role.Color,
				Name:     role.Name,
				Position: role.Position,
				ServerID: uint(guildFlake),
			},
		)
	}

}
func guildDelete(s *discordgo.Session, e *discordgo.GuildDelete) {
	guildFlake, err := snowflake.ParseString(e.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	db.Delete(&Server{Model: gorm.Model{ID: uint(guildFlake)}})

}

func guildUpdate(s *discordgo.Session, e *discordgo.GuildUpdate) {
	guildFlake, err := snowflake.ParseString(e.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}

	roles := make([]Role, len(e.Roles))
	for i, role := range e.Roles {

		roleFlake, err := snowflake.ParseString(role.ID)
		if err != nil {
			log.Fatalf("parse role snowflake: %+v", err)
		}
		roles[i] = Role{
			Model:    gorm.Model{ID: uint(roleFlake)},
			Color:    role.Color,
			Name:     role.Name,
			Position: role.Position,
			ServerID: uint(guildFlake),
		}
	}
	db.Clauses(
		clause.OnConflict{
			UpdateAll: true,
		},
	).CreateInBatches(&roles, len(roles))
}

func guildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guildFlake, err := snowflake.ParseString(m.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	memberFlake, err := snowflake.ParseString(m.Member.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}

	var server Server
	db.First(&server, uint(guildFlake))

	user := User{
		Model:     gorm.Model{ID: uint(memberFlake)},
		InaraUser: InaraUser{},
		Servers:   []*Server{&server},
		Roles:     nil,
	}

	db.FirstOrCreate(&user)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}

	memb, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		log.Fatalf("get member: %+v", err)
	}

	memberFlake, err := snowflake.ParseString(memb.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v, ", err)
	}
	guildFlake, err := snowflake.ParseString(m.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v, ", err)
	}
	user := User{Model: gorm.Model{ID: uint(memberFlake)}}

	roles := make([]Role, len(m.Member.Roles))

	for i, role := range m.Member.Roles {
		roleFlake, err := snowflake.ParseString(role)
		if err != nil {
			log.Fatalf("parse role snowflake: %+v, ", err)
		}
		roles[i] = Role{Model: gorm.Model{ID: uint(roleFlake)}}
	}

	db.Model(&user).Association("Roles").Append(roles)

	var userRoles []Role
	db.Model(&user).Association("Roles").Find(&userRoles)

	db.First(&user)

	var inaraRanks []Role
	inaraRank := Role{InaraRank: inara.Outsider}
	var eliteRanks []Role
	eliteRank := Role{EliteRank: elite.None}
	for _, role := range userRoles {
		if role.InaraRank > inara.Outsider {
			if role.ServerID == uint(guildFlake) {
				inaraRanks = append(inaraRanks, role)
			}

		}
		if role.EliteRank > elite.None {
			if role.ServerID == uint(guildFlake) {
				eliteRanks = append(eliteRanks, role)
			}

		}
	}

	if len(inaraRanks) > 1 {
		var ranksSorted []Role
		From(inaraRanks).Sort(
			func(i, j interface{}) bool {
				i2 := i.(Role)
				j2 := j.(Role)
				return i2.InaraRank < j2.InaraRank
			},
		).ToSlice(&ranksSorted)
		inaraRank = ranksSorted[len(ranksSorted)-1]
		ranksSorted = ranksSorted[:len(ranksSorted)-1]
		for _, role := range ranksSorted {
			err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, snowflake.ID(role.ID).String())
			if err != nil {
				log.Fatalf("remove redundant roles: %+v", err)
			}
		}
	} else if len(inaraRanks) == 1 {
		inaraRank = inaraRanks[0]
	}

	if len(eliteRanks) > 1 {
		var ranksSorted []Role
		From(eliteRanks).Sort(
			func(i, j interface{}) bool {
				i2 := i.(Role)
				j2 := j.(Role)
				return i2.EliteRank < j2.EliteRank
			},
		).ToSlice(&ranksSorted)
		eliteRank = ranksSorted[len(ranksSorted)-1]
		ranksSorted = ranksSorted[:len(ranksSorted)-1]
		for _, role := range ranksSorted {
			err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, snowflake.ID(role.ID).String())
			if err != nil {
				log.Fatalf("remove redundant roles: %+v", err)
			}
		}
	} else if len(eliteRanks) == 1 {
		eliteRank = eliteRanks[0]
	}

	if inaraRank.InaraRank != inara.Outsider {
		if inaraRank.InaraRank != user.InaraRank {
			err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, snowflake.ID(inaraRank.ID).String())
			if err != nil {
				log.Fatalf("remove wrong role: %+v", err)
			}
		}
	}
	if eliteRank.EliteRank != elite.None {
		if eliteRank.EliteRank != user.EliteRank {
			err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, snowflake.ID(eliteRank.ID).String())
			if err != nil {
				log.Fatalf("remove wrong role: %+v", err)
			}
		}
	}

	var rightInaraRole Role
	db.Where("inara_rank = ? AND server_id = ?", user.InaraRank, uint(guildFlake)).Limit(1).Find(&rightInaraRole)

	var rightEliteRole Role
	db.Where("elite_rank = ? AND server_id = ?", user.EliteRank, uint(guildFlake)).Limit(1).Find(&rightEliteRole)

	if rightInaraRole.ID == rightEliteRole.ID {
		if rightInaraRole.ID != 0 && rightEliteRole.ID != 0 {
			err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, snowflake.ID(rightEliteRole.ID).String())
			if err != nil {
				log.Fatalf("add role: %+v", err)
			}
		}

	} else {
		if rightEliteRole.ID != 0 {
			err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, snowflake.ID(rightEliteRole.ID).String())
			if err != nil {
				log.Fatalf("add role: %+v", err)
			}
		}

		if rightInaraRole.ID != 0 {
			err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, snowflake.ID(rightInaraRole.ID).String())
			if err != nil {
				log.Fatalf("add role: %+v", err)
			}
		}

	}

	log.Println(inaraRanks)
	log.Println(eliteRanks)

	log.Println("roles")

	// if strings.Contains(m.Content, "coriolis.io") || strings.Contains(m.Content, "orbis.zone") {
	// 	err := ShipBuildCoriolis(context.TODO(), m.Content, s, m)
	// 	if err != nil {
	// 		log.Fatalf("error on Command Coriolis Ship Build: %+v", err)
	// 	}
	// }
	if strings.Contains(m.Content, "edsy") {
		err := ShipBuildEDSY(context.TODO(), m.Content, s, m)
		if err != nil {
			log.Fatalf("error on Command EDSY Ship Build: %+v", err)
		}
	}

}
