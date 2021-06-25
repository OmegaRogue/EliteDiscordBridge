package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"edDiscord/inara"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

var (
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	InaraKey     = os.Getenv("INARA_KEY")
	EDSMUser     = os.Getenv("EDSM_USER")
	EDSMKey      = os.Getenv("EDSM_KEY")
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Stamp,
	}).With().Timestamp().Caller().Logger()
}

type DiscordContext struct {
	context.Context
	Session *discordgo.Session
	Event   interface{}
}

var db *gorm.DB

var inaraClient *inara.API

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	levelString := flag.String("log-level", zerolog.InfoLevel.String(), "sets log level")
	flag.Parse()
	level, err := zerolog.ParseLevel(*levelString)
	if err != nil {
		log.Err(err).Str("logLevel", *levelString).Msg("parse log level")
	} else {
		log.Trace().Str("logLevel", *levelString).Msg("parse log level")
	}

	zerolog.SetGlobalLevel(level)

	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Err(err).Msg("Initialize Discord Bot")
	}
	log.Trace().Msg("Initialize Discord Bot")
	db, err = gorm.Open(sqlite.Open("eliteDiscord.db"), &gorm.Config{})
	if err != nil {
		log.Err(err).Msg("Open DB")
	}
	log.Trace().Msg("Open DB")
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Err(err).Msg("migrate Users")
	}
	log.Trace().Msg("migrate Users")
	err = db.AutoMigrate(&Server{})
	if err != nil {
		log.Err(err).Msg("migrate Servers")
	}
	log.Trace().Msg("migrate Servers")
	err = db.AutoMigrate(&Role{})
	if err != nil {
		log.Err(err).Msg("migrate Roles")
	}
	log.Trace().Msg("migrate Roles")
	err = db.AutoMigrate(&InaraUser{})
	if err != nil {
		log.Err(err).Msg("migrate InaraUsers")
	}
	log.Trace().Msg("migrate InaraUsers")
	err = db.AutoMigrate(&InaraWing{})
	if err != nil {
		log.Err(err).Msg("migrate InaraWings")
	}
	log.Trace().Msg("migrate InaraWings")
	err = db.AutoMigrate(&InaraSquadron{})
	if err != nil {
		log.Err(err).Msg("migrate InaraSquadrons")
	}
	log.Trace().Msg("migrate InaraSquadrons")

	inaraHead := inara.Header{
		AppName:          "EliteDiscordBridge",
		AppVersion:       "0.0.0",
		IsBeingDeveloped: true,
		APIkey:           InaraKey,
	}
	inaraClient = inara.NewAPI(inaraHead)
	log.Trace().Msg("inizialize Inara API")

	dg.AddHandler(guildCreated)
	log.Trace().Msg("register handler for GuildCreate Event")
	dg.AddHandler(guildDelete)
	log.Trace().Msg("register handler for GuildDelete Event")
	dg.AddHandler(guildUpdate)
	log.Trace().Msg("register handler for GuildUpdate Event")
	dg.AddHandler(guildMemberAdd)
	log.Trace().Msg("register handler for GuildMemberAdd Event")
	dg.AddHandler(messageCreate)
	log.Trace().Msg("register handler for MessageCreate Event")
	dg.AddHandler(roleCreated)
	log.Trace().Msg("register handler for GuildRoleCreate Event")
	dg.AddHandler(roleUpdated)
	log.Trace().Msg("register handler for GuildRoleUpdate Event")

	dg.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlers[i.Data.Name]; ok {
				h(s, i)
			}
		},
	)
	log.Trace().Msg("add handler for InteractionCreate Event")

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Err(err).Msg("open websocket to discord")
	}
	log.Trace().Msg("open websocket to discord")

	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			log.Err(err).Msg("close discord session")
		}
		log.Trace().Msg("close discord session")
	}(dg)

	idle := 42206092983000
	err = dg.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			IdleSince: &idle, Activities: []*discordgo.Activity{
				{
					Name: "Connecting Elite Dangerous APIs since 3307",
					Type: discordgo.ActivityTypeGame,
					URL:  "https://github.com/OmegaRogue/EliteDiscordBridge",
				},
			}, AFK: true, Status: "online",
		},
	)
	if err != nil {
		log.Err(err).Msg("update discord bot status")
	}
	log.Trace().Msg("update discord bot status")

	for i, v := range commands {

		v2, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
		if err != nil {
			log.Err(err).Interface("command", v2).Msg("create slash command")
		}
		commands[i] = v2

		log.Trace().Interface("command", v2).Msg("create slash command")
	}

	log.Info().Msg("Bot is now running. Press CTRL-C to exit.")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for _, v := range commands {
		err = dg.ApplicationCommandDelete(dg.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Err(err).Interface("command", v).Msg("delete slash command")
		}
		log.Trace().Interface("command", v).Msg("delete slash command")
	}

}

var GuildID = "627059611219918869"
