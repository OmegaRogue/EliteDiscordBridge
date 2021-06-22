package main

import (
	"fmt"
	"strconv"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"

	"edDiscord/elite"
	"edDiscord/inara"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/mattn/go-colorable"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(
		&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"component", "category"},
		},
	)
	log.SetOutput(colorable.NewColorableStdout())
}

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
			inaraProfileData, err := inaraClient.GetProfile(i.Data.Options[1].StringValue())
			log.WithField("Command", "Register")
			if err != nil {
				log.Fatalf("get inara profile: %+v", err)
			}

			guildFlake, err := snowflake.ParseString(i.GuildID)
			if err != nil {
				log.WithField("GuildID", i.GuildID).
					Fatalf("parse guild snowflake: %+v", err)
			}
			memberFlake, err := snowflake.ParseString(i.Data.Options[0].UserValue(s).ID)
			if err != nil {
				log.WithField("UserValue", i.Data.Options[0].UserValue(s)).
					Fatalf("parse member snowflake: %+v, ", err)
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

			inaraUser := InaraUser{
				Model:           gorm.Model{ID: uint(inaraProfileData.UserID)},
				UserID:          uint(memberFlake),
				Name:            inaraProfileData.UserName,
				CommanderName:   inaraProfileData.CommanderName,
				Allegiance:      elite.ParseAllegiance(inaraProfileData.PreferredAllegianceName),
				Power:           elite.ParsePower(inaraProfileData.PreferredPowerName),
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

			s.InteractionRespond(
				i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Content: "Hey there! Congratulations, you just registered to EliteDiscordBridge",
					},
				},
			)
		},
		"link": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// As you can see, the name of subcommand (nested, top-level) or subcommand group
			// is provided through arguments.
			switch i.Data.Options[0].Name {
			case "rank":
				role := i.Data.Options[0].Options[0].Options[0].RoleValue(s, i.GuildID)
				roleFlake, err := snowflake.ParseString(role.ID)
				if err != nil {
					log.Fatalf("parse roleFlake: %+v", err)
				}
				dbRole := Role{
					Model: gorm.Model{ID: uint(roleFlake)},
				}

				switch i.Data.Options[0].Options[0].Name {
				case "inara":
					rank := inara.Rank(i.Data.Options[0].Options[0].Options[1].IntValue())
					db.Model(&dbRole).Update("inara_rank", rank)

					s.InteractionRespond(
						i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionApplicationCommandResponseData{
								Content: fmt.Sprintf(
									"Linked Role %s to inara rank %v",
									role.Mention(),
									dbRole.InaraRank.ToString(),
								),
							},
						},
					)
				case "elite":
					rank := elite.Rank(i.Data.Options[0].Options[0].Options[1].IntValue())
					db.Model(&dbRole).Update("elite_rank", rank)
					s.InteractionRespond(
						i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionApplicationCommandResponseData{
								Content: fmt.Sprintf("Linked Role %s to elite dangerous rank %v", role.Mention(), dbRole.EliteRank.ToString()),
							},
						},
					)
				}

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

func AssignInaraRanks(parts []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := snowflake.ParseString(m.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	_, err = snowflake.ParseString(m.Author.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}
	return nil
}
