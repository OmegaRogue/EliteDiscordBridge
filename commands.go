package main

import (
	"fmt"
	"strconv"

	"edDiscord/inara"
	elite "github.com/OmegaRogue/eliteJournal"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "register",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "username",
					Description: "Inara Username/Elite Dangerous Commander Name",
					Required:    true,
				},
			},
		},
		{
			Name: "embed",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Embeds an EDSY Link",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "EDSY URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "link",
			Description: "Setup links between Discord and Elite Dangerous",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "rank",
					Description: "Link Roles to Inara and ED Ranks",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "inara",
							Description: "Link Discord Roles to Inara Ranks",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionRole,
									Name:        "role",
									Description: "Role to link",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionInteger,
									Name:        "role-id",
									Description: "RoleID to link to",
									Required:    true,
									Choices: []*discordgo.ApplicationCommandOptionChoice{
										{
											Name:  "None",
											Value: -1,
										},
										{
											Name:  "Recruit",
											Value: 0,
										},
										{
											Name:  "Reserve",
											Value: 1,
										},
										{
											Name:  "Co-Pilot",
											Value: 2,
										},
										{
											Name:  "Pilot",
											Value: 3,
										},
										{
											Name:  "Wingman",
											Value: 4,
										},
										{
											Name:  "Senior Wingman",
											Value: 5,
										},
										{
											Name:  "Veteran",
											Value: 6,
										},
										{
											Name:  "Flight Leader",
											Value: 7,
										},
										{
											Name:  "Flight Group Leader",
											Value: 8,
										},
										{
											Name:  "Operations Officer",
											Value: 9,
										},
										{
											Name:  "Chief of Staff",
											Value: 10,
										},
										{
											Name:  "Deputy Squadron Commander",
											Value: 11,
										},
										{
											Name:  "Squadron Commander",
											Value: 12,
										},
									},
								},
							},
						},
						{
							Name:        "elite",
							Description: "Link Discord Roles to ED Ranks",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Type:        discordgo.ApplicationCommandOptionRole,
									Name:        "role",
									Description: "Role to link",
									Required:    true,
								},
								{
									Type:        discordgo.ApplicationCommandOptionInteger,
									Name:        "role-id",
									Description: "RoleID to link to",
									Required:    true,
									Choices: []*discordgo.ApplicationCommandOptionChoice{
										{
											Name:  "None",
											Value: -1,
										},
										{
											Name:  "Rookie",
											Value: 0,
										},
										{
											Name:  "Agent",
											Value: 1,
										},
										{
											Name:  "Officer",
											Value: 2,
										},
										{
											Name:  "SeniorOfficer",
											Value: 3,
										},
										{
											Name:  "Leader",
											Value: 4,
										},
									},
								},
							},
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			inaraProfileData, err := inaraClient.GetProfile(i.ApplicationCommandData().Options[1].StringValue())
			if err != nil {
				log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("get inara profile")
				return
			}

			guildFlake, err := snowflake.ParseString(i.GuildID)
			if err != nil {
				log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("parse guild snowflake")
				return
			}
			log.Trace().Caller().Int64("guild", int64(guildFlake)).Msg("parse guild snowflake")
			memberFlake, err := snowflake.ParseString(i.ApplicationCommandData().Options[0].UserValue(s).ID)
			if err != nil {
				log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Interface(
					"UserValue",
					i.ApplicationCommandData().Options[0].UserValue(s),
				).Msg("parse member snowflake")
				return
			}

			var server Server
			db.First(&server, uint(guildFlake))
			user := User{
				Model:   gorm.Model{ID: uint(memberFlake)},
				Servers: []*Server{&server},
			}

			db.FirstOrCreate(&user)

			// img := imgbase64.FromRemote(inaraProfileData.AvatarImageURL)

			inaraSquadron := InaraSquadron{
				Model: gorm.Model{ID: uint(inaraProfileData.CommanderSquadron.SquadronID)},
				Name:  inaraProfileData.CommanderSquadron.SquadronName,
				URL:   inaraProfileData.CommanderSquadron.InaraURL,
			}

			inaraWing := InaraWing{
				Model:           gorm.Model{ID: uint(inaraProfileData.CommanderWing.WingID)},
				Name:            inaraProfileData.CommanderWing.WingName,
				URL:             inaraProfileData.CommanderWing.InaraURL,
				InaraSquadronID: inaraSquadron.ID,
			}

			allegiance, err := elite.ParseAllegiance(inaraProfileData.PreferredAllegianceName)
			if err != nil {
				log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("parse allegiance")
				return
			}
			power, err := elite.ParsePower(inaraProfileData.PreferredPowerName)
			if err != nil {
				log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("parse power")
				return
			}
			inaraUser := InaraUser{
				Model:           gorm.Model{ID: uint(inaraProfileData.UserID)},
				UserID:          uint(memberFlake),
				Name:            inaraProfileData.UserName,
				CommanderName:   inaraProfileData.CommanderName,
				Allegiance:      allegiance,
				Power:           power,
				AvatarURL:       inaraProfileData.AvatarImageURL,
				URL:             inaraProfileData.InaraURL,
				InaraWingID:     inaraWing.ID,
				InaraSquadronID: inaraSquadron.ID,
				CombatRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("combat")).
					Select(getPilotRankValue).
					First().(int),
				CombatProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("combat")).
					Select(getPilotRankProgress).
					First().(float64),
				TradeRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("trade")).
					Select(getPilotRankValue).
					First().(int),
				TradeProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("trade")).
					Select(getPilotRankProgress).
					First().(float64),
				ExplorationRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("exploration")).
					Select(getPilotRankValue).
					First().(int),
				ExplorationProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("exploration")).
					Select(getPilotRankProgress).
					First().(float64),
				CqcRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("cqc")).
					Select(getPilotRankValue).
					First().(int),
				CqcProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("cqc")).
					Select(getPilotRankProgress).
					First().(float64),
				SoldierRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("soldier")).
					Select(getPilotRankValue).
					First().(int),
				SoldierProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("soldier")).
					Select(getPilotRankProgress).
					First().(float64),
				ExobiologistRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("exobiologist")).
					Select(getPilotRankValue).
					First().(int),
				ExobiologistProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("exobiologist")).
					Select(getPilotRankProgress).
					First().(float64),
				EmpireRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("empire")).
					Select(getPilotRankValue).
					First().(int),
				EmpireProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("empire")).
					Select(getPilotRankProgress).
					First().(float64),
				FederationRank: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("federation")).
					Select(getPilotRankValue).
					First().(int),
				FederationProgress: From(inaraProfileData.CommanderRanksPilot).
					Where(getPilotRankPredicate("federation")).
					Select(getPilotRankProgress).
					First().(float64),
			}

			db.Clauses(
				clause.OnConflict{
					UpdateAll: true,
				},
			).Create(&inaraSquadron)

			db.Clauses(
				clause.OnConflict{
					UpdateAll: true,
				},
			).Create(&inaraWing)

			db.Clauses(
				clause.OnConflict{
					UpdateAll: true,
				},
			).Create(&inaraUser)

			err = s.InteractionRespond(
				i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Hey there! Congratulations, you just registered to EliteDiscordBridge",
					},
				},
			)
			if err != nil {
				log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("respond")
				return
			}
		},
		"link": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// As you can see, the name of subcommand (nested, top-level) or subcommand group
			// is provided through arguments.
			switch i.ApplicationCommandData().Options[0].Name {
			case "rank":
				role := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID)
				roleFlake, err := snowflake.ParseString(role.ID)
				if err != nil {
					log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("parse roleFlake")
					return
				}
				log.Trace().Caller().Interface("role", role).Msg("parse role snowflake")
				dbRole := Role{
					Model: gorm.Model{ID: uint(roleFlake)},
				}

				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "inara":
					rank := inara.Rank(i.ApplicationCommandData().Options[0].Options[0].Options[1].IntValue())
					db.Model(&dbRole).Update("inara_rank", rank)

					err := s.InteractionRespond(
						i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf(
									"Linked Role %s to inara rank %v",
									role.Mention(),
									dbRole.InaraRank,
								),
							},
						},
					)
					if err != nil {
						log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("error on respond")
						return
					}
				case "elite":
					rank := elite.Rank(i.ApplicationCommandData().Options[0].Options[0].Options[1].IntValue())
					db.Model(&dbRole).Update("elite_rank", rank)
					err := s.InteractionRespond(
						i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf(
									"Linked Role %s to elite dangerous rank %v",
									role.Mention(),
									dbRole.EliteRank,
								),
							},
						},
					)
					if err != nil {
						log.Err(err).Stack().Caller().Interface("InteractionCreate", i).Msg("error on respond")
						return
					}
				}

			}

		},

		"embed": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// As you can see, the name of subcommand (nested, top-level) or subcommand group
			// is provided through arguments.
			err := ShipBuildEDSY(browserContext, i.ApplicationCommandData().Options[0].StringValue(), s, nil, i)
			if err != nil {
				log.Err(err).Stack().Caller().Interface("interaction", i).Msg("Command EDSY Ship Build")
				return
			}

		},
	}
)

func getPilotRankValue(i interface{}) interface{} {
	rank, _ := strconv.Atoi(i.(inara.PilotRank).RankValue)
	return rank
}
func getPilotRankProgress(i interface{}) interface{} {
	return i.(inara.PilotRank).RankProgress
}

func getPilotRankPredicate(rankName string) func(i interface{}) bool {
	return func(i interface{}) bool {
		return i.(inara.PilotRank).RankName == rankName
	}
}
