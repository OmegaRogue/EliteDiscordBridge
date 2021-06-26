package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	gormLog2 "edDiscord/gormLog"
	"edDiscord/inara"
	elite "github.com/OmegaRogue/eliteJournal"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
	"github.com/halink0803/zerolog-graylog-hook/graylog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	InaraKey     = os.Getenv("INARA_KEY")
	EDSMUser     = os.Getenv("EDSM_USER")
	EDSMKey      = os.Getenv("EDSM_KEY")
)

func init() {

}

type DiscordContext struct {
	context.Context
	Session *discordgo.Session
	Event   interface{}
}

var db *gorm.DB

var inaraClient *inara.API

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var logMap = map[int]zerolog.Level{
	discordgo.LogError:         zerolog.ErrorLevel,
	discordgo.LogWarning:       zerolog.WarnLevel,
	discordgo.LogInformational: zerolog.InfoLevel,
	discordgo.LogDebug:         zerolog.DebugLevel,
}

var logMap2 = map[zerolog.Level]int{
	zerolog.DebugLevel: discordgo.LogDebug,
	zerolog.InfoLevel:  discordgo.LogInformational,
	zerolog.WarnLevel:  discordgo.LogWarning,
	zerolog.ErrorLevel: discordgo.LogError,
	zerolog.PanicLevel: -1,
	zerolog.FatalLevel: -1,
	zerolog.TraceLevel: discordgo.LogDebug,
	zerolog.NoLevel:    -1,
	zerolog.Disabled:   -1,
}

var logMap3 = map[zerolog.Level]gormlogger.LogLevel{
	zerolog.DebugLevel: gormlogger.Info,
	zerolog.InfoLevel:  gormlogger.Info,
	zerolog.WarnLevel:  gormlogger.Warn,
	zerolog.ErrorLevel: gormlogger.Error,
	zerolog.PanicLevel: gormlogger.Silent,
	zerolog.FatalLevel: gormlogger.Silent,
	zerolog.TraceLevel: gormlogger.Info,
	zerolog.NoLevel:    gormlogger.Silent,
	zerolog.Disabled:   gormlogger.Silent,
}

var discordOpCode = map[int64]string{
	0:  "Dispatch",
	1:  "Heartbeat",
	2:  "Identify",
	3:  "Presence Update",
	4:  "Voice State Update",
	6:  "Resume",
	7:  "Reconnect",
	8:  "Request Guild Members",
	9:  "Invalid Session",
	10: "Hello",
	11: "Heartbeat ACK",
}

func AnyInt(n interface{}) interface{} {
	switch n := n.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case uint:
		return uint64(n)
	case uintptr:
		return uint64(n)
	case uint8:
		return uint64(n)
	case uint16:
		return uint64(n)
	case uint32:
		return uint64(n)
	case uint64:
		return uint64(n)
	}
	return nil
}

func AnySInt(n interface{}) (int64, error) {
	switch n := n.(type) {
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	}
	return 0, errors.New("n is not a signed int")
}

const ComponentFieldName = "component"

var SyncTimer *time.Timer
var SyncInterval time.Duration

func main() {
	SyncInterval = time.Hour

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		files := strings.Split(file, "/")
		file = files[len(files)-1]
		return file + ":" + strconv.Itoa(line)
	}
	// log.Logger = log.Output(horizontal.ConsoleWriter{Out: os.Stdout}).
	// 	With().Timestamp().Str(
	// 	ComponentFieldName,
	// 	"Bot",
	// ).Logger()
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.Stamp,
			PartsOrder: []string{
				zerolog.TimestampFieldName,
				zerolog.LevelFieldName,
				zerolog.CallerFieldName,
				ComponentFieldName,
				zerolog.MessageFieldName,
			},
		},
	).With().Timestamp().Str(ComponentFieldName, "Bot").Logger()

	hook, err := graylog.NewGraylogHook("udp://127.0.0.1:12201")
	if err != nil {
		panic(err)
	}
	//Set global logger with graylog hook
	log.Logger = log.Hook(hook).With().Timestamp().Str(ComponentFieldName, "Bot").Logger()
	levelString := flag.String("log-level", zerolog.InfoLevel.String(), "sets log level")
	flag.Parse()
	level, err := zerolog.ParseLevel(*levelString)
	if err != nil {
		log.Err(err).Stack().Caller().Str("logLevel", *levelString).Msg("parse log level")
	} else {
		log.Info().Caller().Str("logLevel", *levelString).Msg("parse log level")
	}

	zerolog.SetGlobalLevel(level)

	discordLog := log.With().Str(ComponentFieldName, "Discord").Logger()
	discordgo.Logger = DiscordLogParse(discordLog)

	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Err(err).Stack().Caller().Msg("Initialize Discord Bot")
	} else {
		log.Debug().Caller().Msg("Initialize Discord Bot")
	}
	dg.LogLevel = logMap2[level]

	gormLog := gormLog2.NewWithLogger(log.With().Str(ComponentFieldName, "GORM").Logger())
	gormLog.SourceField = "caller"
	db, err = gorm.Open(
		sqlite.Open("eliteDiscord.db"), &gorm.Config{
			//Logger: gormlogger.Default.LogMode(logMap3[level]),
			//Logger: gormLog,
			Logger: gormLog,
		},
	)
	if err != nil {
		log.Err(err).Stack().Caller().Msg("Open DB")
	} else {
		log.Info().Caller().Msg("Open DB")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Err(err).Stack().Caller().Msg("migrate Users")
	} else {
		log.Debug().Caller().Msg("migrate Users")
	}
	err = db.AutoMigrate(&Server{})
	if err != nil {
		log.Err(err).Stack().Caller().Msg("migrate Servers")
	} else {
		log.Debug().Caller().Msg("migrate Servers")
	}
	err = db.AutoMigrate(&Role{})
	if err != nil {
		log.Err(err).Stack().Caller().Msg("migrate Roles")
	} else {
		log.Debug().Caller().Msg("migrate Roles")
	}
	err = db.AutoMigrate(&InaraUser{})
	if err != nil {
		log.Err(err).Stack().Caller().Msg("migrate InaraUsers")
	} else {
		log.Debug().Caller().Msg("migrate InaraUsers")
	}
	err = db.AutoMigrate(&InaraWing{})
	if err != nil {
		log.Err(err).Stack().Caller().Msg("migrate InaraWings")
	} else {
		log.Debug().Caller().Msg("migrate InaraWings")
	}
	err = db.AutoMigrate(&InaraSquadron{})
	if err != nil {
		log.Err(err).Stack().Caller().Msg("migrate InaraSquadrons")
	} else {
		log.Debug().Caller().Msg("migrate InaraSquadrons")
	}
	inaraHead := inara.Header{
		AppName:          "EliteDiscordBridge",
		AppVersion:       "0.0.0",
		IsBeingDeveloped: true,
		APIkey:           InaraKey,
	}
	inaraClient = inara.NewAPI(inaraHead)
	log.Info().Caller().Msg("inizialize Inara API")

	dg.AddHandler(guildCreated)
	log.Debug().Caller().Msg("register handler for GuildCreate Event")
	dg.AddHandler(guildDelete)
	log.Debug().Caller().Msg("register handler for GuildDelete Event")
	dg.AddHandler(guildUpdate)
	log.Debug().Caller().Msg("register handler for GuildUpdate Event")
	dg.AddHandler(guildMemberAdd)
	log.Debug().Caller().Msg("register handler for GuildMemberAdd Event")
	dg.AddHandler(messageCreate)
	log.Debug().Caller().Msg("register handler for MessageCreate Event")
	dg.AddHandler(roleCreated)
	log.Debug().Caller().Msg("register handler for GuildRoleCreate Event")
	dg.AddHandler(roleUpdated)
	log.Debug().Caller().Msg("register handler for GuildRoleUpdate Event")
	// wsapi.go:597:onEvent()
	// ../../../../../../home/omegarogue/go/pkg/mod/github.com/bwmarrin/discordgo@v0.23.3-0.20210617211910-e72c457cb4ae/logging.go:77
	dg.AddHandler(
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlers[i.Data.Name]; ok {
				h(s, i)
			}
		},
	)
	log.Debug().Caller().Msg("register handler for InteractionCreate Event")

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Err(err).Stack().Caller().Msg("open websocket to discord")
	}
	log.Info().Caller().Msg("open websocket to discord")

	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			log.Err(err).Stack().Caller().Msg("close discord session")
		}
		log.Info().Caller().Msg("close discord session")
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
		log.Err(err).Stack().Caller().Msg("update discord bot status")
	}
	log.Info().Caller().Msg("update discord bot status")

	for i, v := range commands {

		v2, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
		if err != nil {
			log.Err(err).Stack().Caller().Interface("command", v2).Int("i", i).Msg("create slash command")
			continue
		}
		commands[i] = v2

	}

	log.Info().Caller().Msg("is now running. Press CTRL-C to exit.")

	SyncTimer = time.NewTimer(SyncInterval)

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	go SyncRoutine(sc)
	<-sc

	for i, v := range commands {
		err = dg.ApplicationCommandDelete(dg.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Err(err).Stack().Caller().Int("i", i).Interface("command", v).Msg("delete slash command")
			continue
		}
	}

}

func SyncRoutine(sc chan os.Signal) {
	log.Info().Caller().Msg("Heartbeat started")
	select {
	case <-SyncTimer.C:
		defer SyncTimer.Reset(SyncInterval)
		defer func() { go SyncRoutine(sc) }()
		log.Info().Caller().Msg("Heartbeat")
		InaraSync()
	case <-sc:
		return

	}

}

func InaraSync() {
	var users []InaraUser
	db.Find(&users)
	var userNames []string
	From(users).Select(
		func(i interface{}) interface{} {
			return i.(InaraUser).Name
		}).ToSlice(&userNames)
	profiles, err := inaraClient.GetProfiles(userNames)
	if err != nil {
		log.Err(err).Stack().Caller().Interface("profiles", profiles).Msg("get profiles")
	}
	log.Info().Stack().Caller().Interface("profiles", profiles).Msg("get profiles")

	for i, inaraProfileData := range profiles {
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
			UserID:          users[i].UserID,
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

		db.Save(&inaraSquadron)
		db.Save(&inaraWing)
		db.Save(&inaraUser)

	}

}

func DiscordLogParse(discordLog zerolog.Logger) func(msgL int, caller int, format string, a ...interface{}) {
	return func(msgL, caller int, format string, a ...interface{}) {

		e := discordLog.WithLevel(logMap[msgL]).CallerSkipFrame(1).Caller(caller)

		if opI := strings.Index(format, "Op: %d"); opI != -1 {
			pC := strings.Count(format[:opI], "%")

			e.Str("op", discordOpCode[int64(a[pC].(int))])
			From(a).Except(From(a).Where(func(i interface{}) bool { return i == a[pC] })).ToSlice(&a)
			format = strings.Join(strings.Split(format, "Op: %d"), "")

		}
		if seqI := strings.Index(format, "Seq: %d"); seqI != -1 {
			pC := strings.Count(format[:seqI], "%")

			e.Interface("seq", a[pC])
			From(a).Except(From(a).Where(func(i interface{}) bool { return i == a[pC] })).ToSlice(&a)
			format = strings.Join(strings.Split(format, "Seq: %d"), "")

		}
		if seqI := strings.Index(format, "seq %d"); seqI != -1 {
			pC := strings.Count(format[:seqI], "%")

			e.Interface("seq", a[pC])
			From(a).Except(From(a).Where(func(i interface{}) bool { return i == a[pC] })).ToSlice(&a)
			format = strings.Join(strings.Split(format, "seq %d"), "")

		}
		if typeI := strings.Index(format, "Type: %s"); typeI != -1 {
			pC := strings.Count(format[:typeI], "%")

			data, hasData := a[pC].(string)
			if hasData {
				if data != "" {
					e.Str("type", data)
				}
				From(a).Except(From(a).Where(func(i interface{}) bool { return i == a[pC] })).ToSlice(&a)
				format = strings.Join(strings.Split(format, "Type: %s"), "")
			}
		}
		if dataI := strings.Index(format, "Data: %s"); dataI != -1 {
			pC := strings.Count(format[:dataI], "%")

			data, hasData := a[pC].(string)
			if hasData {
				if data != "null" {
					// data, _ = strconv.Unquote(data)

					var js map[string]interface{}
					json.Unmarshal([]byte(data), &js)
					trace, ok := js["_trace"]
					if ok {
						tracer := trace.([]interface{})
						var traceBack []interface{}
						for _, i2 := range tracer {
							var js3 []interface{}
							json.Unmarshal([]byte(i2.(string)), &js3)
							traceBack = append(traceBack, js3)

						}
						js["_trace"] = traceBack
						data2, _ := json.Marshal(js)
						data = string(data2)
					}

					e.RawJSON("data", []byte(data))
				}
				From(a).Except(From(a).Where(func(i interface{}) bool { return i == a[pC] })).ToSlice(&a)
				format = strings.Join(strings.Split(format, "Data: %s"), "")
			}

		}

		msg := strings.Trim(format, ", \n")
		if msg == "" {
			e.Send()
		} else {
			e.Msgf(msg, a...)
		}

	}
}

var GuildID = "627059611219918869"
