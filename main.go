package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/go-resty/resty/v2"
	_ "modernc.org/sqlite"
)

var (
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	InaraKey     = os.Getenv("INARA_KEY")
	EDSMUser     = os.Getenv("EDSM_USER")
	EDSMKey      = os.Getenv("EDSM_KEY")
)

var db *sql.DB
var inara *resty.Client
var inaraHead InaraHeader = InaraHeader{
	AppName:          "EliteDiscordBridge",
	AppVersion:       "0.0.0",
	IsBeingDeveloped: true,
	APIkey:           InaraKey,
}

func main() {
	ctx := context.TODO()

	dg, err := discordgo.New("Bot " + DiscordToken)
	if err != nil {
		log.Fatalf("error on Initialize Discord bot: %+v", err)
	}
	db, err = sql.Open("sqlite", "eliteDiscord.db")
	if err != nil {
		log.Fatalf("error on open DB: %+v", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		log.Panicf("error on ping db: %+v", err)
	}

	inara = resty.New()

	dg.AddHandler(guildCreated)
	dg.AddHandler(guildDelete)
	dg.AddHandler(guildMemberAdd)
	dg.AddHandler(guildMemberRemove)
	dg.AddHandler(guildMemberUpdate)
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error on open websocket: %+v", err)
	}
	defer dg.Close()

	idle := 42206092983000
	dg.UpdateStatusComplex(discordgo.UpdateStatusData{IdleSince: &idle, Activities: []*discordgo.Activity{{
		Name: "Connecting Elite Dangerous APIs since 3307",
		Type: discordgo.ActivityTypeGame,
		URL:  "",
	}}, AFK: true, Status: "online"})

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
func RegisterMember(ctx context.Context, guild, member string) {
	guildFlake, err := snowflake.ParseString(guild)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	memberFlake, err := snowflake.ParseString(member)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}

	db.ExecContext(ctx,
		"INSERT INTO users (userID) VALUES ($1);",
		memberFlake.Int64(),
	)

	db.ExecContext(ctx,
		"INSERT INTO guildUser (userID, serverID) VALUES ($1, $2);",
		memberFlake.Int64(),
		guildFlake.Int64(),
	)

}
func RemoveMember(ctx context.Context, m *discordgo.Member) {
	guildFlake, err := snowflake.ParseString(m.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	memberFlake, err := snowflake.ParseString(m.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}

	db.ExecContext(ctx,
		"DELETE FROM guildUser WHERE userID = $1 AND serverID = $2;",
		memberFlake.Int64(),
		guildFlake.Int64(),
	)
	CheckIfMemberDelete(ctx, m)

}

func RemoveGuild(ctx context.Context, g *discordgo.Guild) {
	flake, err := snowflake.ParseString(g.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}

	db.ExecContext(context.TODO(),
		"DELETE FROM servers WHERE serverID = $1;",
		flake.Int64(),
	)
	db.ExecContext(context.TODO(),
		"DELETE FROM guildUser WHERE serverID = $1;",
		flake.Int64(),
	)

}

func CheckIfMemberDelete(ctx context.Context, m *discordgo.Member) {
	memberFlake, err := snowflake.ParseString(m.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}
	var n int
	err = db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT serverID) FROM guildUser WHERE userID = $1", memberFlake).Scan(&n)
	if err != nil {
		log.Fatalf("get member servers: %+v", err)
	}
	if n == 0 {
		db.ExecContext(ctx, "DELETE FROM users WHERE userID = $1;", memberFlake.Int64())
	}
}

func MemberCleanup(ctx context.Context) {

	rows, err := db.QueryContext(ctx, "SELECT userID, Count(distinct serverID) FROM guildUser GROUP BY userID;")
	if err != nil {
		log.Fatalf("get member servers: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var n int
		rows.Scan(&id, &n)
		log.Println(n)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func guildCreated(s *discordgo.Session, m *discordgo.GuildCreate) {
	guildFlake, err := snowflake.ParseString(m.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	db.ExecContext(context.TODO(),
		"INSERT INTO servers (serverID) VALUES ($1);",
		guildFlake.Int64(),
	)

	for !s.DataReady {
		log.Println("WAIT")
	}
	for _, member := range m.Members {
		log.Println(member.User.Username)
		RegisterMember(context.TODO(), m.ID, member.User.ID)
	}
}
func guildDelete(s *discordgo.Session, m *discordgo.GuildDelete) {
	flake, err := snowflake.ParseString(m.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}

	db.ExecContext(context.TODO(),
		"DELETE FROM servers WHERE serverID = $1;",
		flake.Int64(),
	)
	db.ExecContext(context.TODO(),
		"DELETE FROM guildUser WHERE serverID = $1;",
		flake.Int64(),
	)
	MemberCleanup(context.TODO())
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

	db.ExecContext(context.TODO(),
		"INSERT INTO users (userID) VALUES ($1);",
		memberFlake.Int64(),
	)

	db.ExecContext(context.TODO(),
		"INSERT INTO guildUser (userID, serverID) VALUES ($1, $2);",
		memberFlake.Int64(),
		guildFlake.Int64(),
	)

}

func guildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	guildFlake, err := snowflake.ParseString(m.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	memberFlake, err := snowflake.ParseString(m.Member.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}
	MemberCleanup(context.TODO())
	db.ExecContext(context.TODO(),
		"DELETE FROM guildUser WHERE userID = $1 AND serverID = $2;",
		memberFlake.Int64(),
		guildFlake.Int64(),
	)
}

func guildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	guildFlake, err := snowflake.ParseString(m.GuildID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	memberFlake, err := snowflake.ParseString(m.Member.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}

	fmt.Sprint(guildFlake, memberFlake)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	RegisterMember(context.TODO(), m.GuildID, m.Author.ID)
	if strings.HasPrefix(m.Content, "!register") {
		parts := strings.Split(m.Content, " ")
		log.Println(parts)
		data := InaraData{
			Header: inaraHead,
			Events: []InaraEvent{
				{
					EventName:      "getCommanderProfile",
					EventTimestamp: time.Now().Format("2006-01-02T15:04:05Z"),
					EventCustomID:  1234,
					EventData: struct {
						SearchName string `json:"searchName"`
					}{SearchName: parts[1]},
				},
			},
		}
		r, _ := inara.R().SetBody(data).Post("https://inara.cz/inapi/v1/")
		fmt.Sprint(r)
		fmt.Println(string(r.Body()))
		var p *InaraResponse
		json.Unmarshal(r.Body(), p)

		memberFlake, err := snowflake.ParseString(m.Author.ID)
		if err != nil {
			log.Fatalf("parse member snowflake: %+v", err)
		}

		db.ExecContext(context.TODO(),
			"UPDATE users SET inaraUserName = $2, eliteDangerousUserName = $3, squadronRank = $4, squadron = $5 WHERE userID = $1;",
			memberFlake.Int64(),
		)

	}

}
