package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
)

func guildCreated(s *discordgo.Session, m *discordgo.GuildCreate) {
	guildFlake, err := snowflake.ParseString(m.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}
	db.ExecContext(
		context.TODO(),
		"INSERT INTO servers (serverID) VALUES ($1);",
		guildFlake.Int64(),
	)
}
func guildDelete(s *discordgo.Session, m *discordgo.GuildDelete) {
	flake, err := snowflake.ParseString(m.ID)
	if err != nil {
		log.Fatalf("parse guild snowflake: %+v", err)
	}

	db.ExecContext(
		context.TODO(),
		"DELETE FROM servers WHERE serverID = $1;",
		flake.Int64(),
	)
	db.ExecContext(
		context.TODO(),
		"DELETE FROM guildUser WHERE serverID = $1;",
		flake.Int64(),
	)
	//MemberCleanup(context.TODO())
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

	db.ExecContext(
		context.TODO(),
		"INSERT INTO users (userID) VALUES ($1);",
		memberFlake.Int64(),
	)

	db.ExecContext(
		context.TODO(),
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
	//MemberCleanup(context.TODO())
	db.ExecContext(
		context.TODO(),
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
		log.Println("register")
		parts := strings.Split(m.Content, " ")
		err := RegisterCommand(context.TODO(), parts, s, m)
		if err != nil {
			log.Fatalf("error on Command Register: %+v", err)
		}
	}
	if strings.HasPrefix(m.Content, "!build") {
		log.Println("build")
		parts := strings.Split(m.Content, " ")
		err := ShipBuildCommand(context.TODO(), parts, s, m)
		if err != nil {
			log.Fatalf("error on Command Ship Build: %+v", err)
		}
	}

}
