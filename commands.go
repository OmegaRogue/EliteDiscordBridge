package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/polds/imgbase64"
)

func RegisterCommand(ctx context.Context, parts []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
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
	r, err := inara.R().SetBody(data).Post("https://inara.cz/inapi/v1/")
	if err != nil {
		return fmt.Errorf("inara getProfile: %w", err)
	}
	var p = new(InaraResponse)
	err = json.Unmarshal(r.Body(), p)
	if err != nil {
		return fmt.Errorf("unmarshal inara profile: %w", err)
	}
	fmt.Println(p)

	memberFlake, err := snowflake.ParseString(m.Author.ID)
	if err != nil {
		return fmt.Errorf("parse member snowflake: %w", err)
	}

	inaraProfileData := p.Events[0].EventData
	_, err = db.ExecContext(
		ctx,
		"UPDATE users SET inaraUserID = $2 WHERE userID = $1;",
		memberFlake.Int64(),
		inaraProfileData.
			UserID,
	)
	if err != nil {
		return fmt.Errorf("update User Inara Data: %w", err)
	}

	_, err = db.ExecContext(
		ctx,
		"INSERT OR IGNORE INTO inaraUsers (id, name, commanderName, allegiance, power, avatarURL, url) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		inaraProfileData.UserID,
		inaraProfileData.UserName,
		inaraProfileData.CommanderName,
		ParseAllegiance(inaraProfileData.PreferredAllegianceName),
		ParsePower(inaraProfileData.PreferredPowerName),
		inaraProfileData.AvatarImageURL,
		inaraProfileData.InaraURL,
	)
	if err != nil {
		return fmt.Errorf("insert inara data: %w", err)
	}
	img := imgbase64.FromRemote(inaraProfileData.AvatarImageURL)
	fmt.Println(img)
	hook, err := s.WebhookCreate(m.ChannelID, "EliteDiscordBridge", "")
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	fmt.Printf("%+v\n", hook)
	_, err = s.WebhookExecute(
		hook.ID, hook.Token, false, &discordgo.WebhookParams{
			Content:   "Registered.",
			Username:  inaraProfileData.UserName,
			AvatarURL: inaraProfileData.AvatarImageURL,
		},
	)
	if err != nil {
		return fmt.Errorf("send webhook message: %w", err)
	}
	return nil
}
