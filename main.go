package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"edDiscord/inara"
	"github.com/bwmarrin/discordgo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	InaraKey     = os.Getenv("INARA_KEY")
	EDSMUser     = os.Getenv("EDSM_USER")
	EDSMKey      = os.Getenv("EDSM_KEY")
)

var db *gorm.DB

var inaraClient *inara.API

func main() {

	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalf("error on Initialize Discord bot: %+v", err)
	}
	db, err = gorm.Open(sqlite.Open("eliteDiscord.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("error on open DB: %+v", err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Server{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&InaraUser{})
	db.AutoMigrate(&InaraWing{})
	db.AutoMigrate(&InaraSquadron{})

	inaraHead := inara.Header{
		AppName:          "EliteDiscordBridge",
		AppVersion:       "0.0.0",
		IsBeingDeveloped: true,
		APIkey:           InaraKey,
	}
	inaraClient = inara.NewAPI(inaraHead)

	dg.AddHandler(guildCreated)
	dg.AddHandler(guildDelete)
	dg.AddHandler(guildUpdate)
	dg.AddHandler(guildMemberAdd)
	dg.AddHandler(messageCreate)
	dg.AddHandler(roleCreated)
	dg.AddHandler(roleUpdated)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.Data.Name]; ok {
			h(s, i)
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Panicf("error on open websocket: %+v", err)
	}
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			log.Panicf("error on close discord session: %+v", err)
		}
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
		log.Panicf("error on update status: %+v", err)
	}

	for i, v := range commands {
		v2, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		commands[i] = v2
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for _, v := range commands {
		err = dg.ApplicationCommandDelete(dg.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

}

var GuildID = "627059611219918869"
