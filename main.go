package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
var inaraHead = InaraHeader{
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
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panicf("error on close db: %+v", err)

		}
	}(db)
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

	_, err = db.ExecContext(
		ctx,
		"INSERT INTO users (userID) VALUES ($1);",
		memberFlake.Int64(),
	)
	if err != nil {
		log.Fatalf("insert user record: %+v", err)
	}

	_, err = db.ExecContext(
		ctx,
		"INSERT INTO guildUser (userID, serverID) VALUES ($1, $2);",
		memberFlake.Int64(),
		guildFlake.Int64(),
	)
	if err != nil {
		log.Fatalf("insert server user record: %+v", err)
	}

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

	_, err = db.ExecContext(
		ctx,
		"DELETE FROM guildUser WHERE userID = $1 AND serverID = $2;",
		memberFlake.Int64(),
		guildFlake.Int64(),
	)
	if err != nil {
		log.Fatalf("remove server user record: %+v", err)
	}
	CheckIfMemberDelete(ctx, m)

}

func RemoveGuild(ctx context.Context, g *discordgo.Guild) {
	flake, err := snowflake.ParseString(g.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}

	_, err = db.ExecContext(
		context.TODO(),
		"DELETE FROM servers WHERE serverID = $1;",
		flake.Int64(),
	)
	if err != nil {
		log.Fatalf("remove server record: %+v", err)
	}
	_, err = db.ExecContext(
		context.TODO(),
		"DELETE FROM guildUser WHERE serverID = $1;",
		flake.Int64(),
	)
	if err != nil {
		log.Fatalf("remove server user records: %+v", err)
	}

}

func CheckIfMemberDelete(ctx context.Context, m *discordgo.Member) {
	memberFlake, err := snowflake.ParseString(m.User.ID)
	if err != nil {
		log.Fatalf("parse member snowflake: %+v", err)
	}
	var n int
	err = db.QueryRowContext(
		ctx,
		"SELECT COUNT(DISTINCT serverID) FROM guildUser WHERE userID = $1",
		memberFlake,
	).Scan(&n)
	if err != nil {
		log.Fatalf("get member servers: %+v", err)
	}
	if n == 0 {
		_, err = db.ExecContext(ctx, "DELETE FROM users WHERE userID = $1;", memberFlake.Int64())
		if err != nil {
			log.Fatalf("remove user record: %+v", err)
		}

	}
}

func MemberCleanup(ctx context.Context) error {

	rows, err := db.QueryContext(ctx, "SELECT userID, Count(distinct serverID) FROM guildUser GROUP BY userID;")
	if err != nil {
		return fmt.Errorf("get member servers: %+v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panicf("close fetched roles: %+v", err)
		}
	}(rows)

	for rows.Next() {
		var id int
		var n int
		err := rows.Scan(&id, &n)
		if err != nil {
			return fmt.Errorf("scan rows: %w", err)
		}
		log.Println(n)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error in rows: %w", err)
	}
	return nil
}
