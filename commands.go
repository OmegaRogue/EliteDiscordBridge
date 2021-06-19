package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/chromedp/chromedp"
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

func ShipBuildCommand(ctx context.Context, parts []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	url := strings.Join(parts[1:], " ")
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var example string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`#build > button:nth-child(7)`),
		// find and click "Example" link
		chromedp.Click(`#build > button:nth-child(7)`, chromedp.NodeVisible),
		// retrieve the text of the textarea
		chromedp.Value(`textarea.cb`, &example),
	)
	if err != nil {
		return fmt.Errorf("browser: %w", err)
	}

	var t ShipLoadout
	t.UnmarshalJSON([]byte(example))

	urlName := strings.ToLower(strings.Replace(strings.Replace(t.Ship, " ", "-", -1), ".", "", -1))

	s.ChannelMessageDelete(m.ChannelID, m.ID)

	utils := "​"
	for _, v := range t.Components.Utility {

		if v != nil {
			v2 := v.(map[string]interface{})
			utils += fmt.Sprintf("**%v%v %v**\n", v2["class"], v2["rating"], v2["group"])
		} else {
			utils += "**EMPTY**\n"
		}

	}
	optional := "​"
	for _, v := range t.Components.Internal {

		if v != nil {
			v2 := v.(map[string]interface{})
			optional += fmt.Sprintf("**%v%v %v**\n", v2["class"], v2["rating"], v2["group"])
		} else {
			optional += "**EMPTY**\n"
		}

	}
	hard := "​"
	for _, v := range t.Components.Hardpoints {

		if v != nil {
			v2 := v.(map[string]interface{})
			hard += fmt.Sprintf("**%v%v %v %v**\n", v2["class"], v2["rating"], v2["mount"], v2["group"])
		} else {
			hard += "**EMPTY**\n"
		}

	}

	_, err = s.ChannelMessageSendEmbed(
		m.ChannelID, &discordgo.MessageEmbed{
			URL:         url,
			Type:        discordgo.EmbedTypeRich,
			Title:       t.Name + "​",
			Description: t.Ship + "​",
			Timestamp:   "",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:    fmt.Sprintf("https://edassets.org/static/img/ship-schematics/qohen-leth/%v.png", urlName),
				Width:  3000,
				Height: 3000,
			},
			Color: 0xC06400,
			Provider: &discordgo.MessageEmbedProvider{
				URL:  "https://coriolis.io",
				Name: "Coriolis",
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name:    m.Author.Username + "​",
				IconURL: m.Author.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Speed",
					Value:  fmt.Sprint(t.Stats.Speed),
					Inline: true,
				},
				{
					Name:   "Boost",
					Value:  fmt.Sprint(t.Stats.Boost),
					Inline: true,
				},
				{
					Name:   "​",
					Value:  "​",
					Inline: true,
				},
				{
					Name:   "Max Jump Range",
					Value:  fmt.Sprint(t.Stats.FullTankRange),
					Inline: true,
				},
				{
					Name:   "Unladen Jump Range",
					Value:  fmt.Sprint(t.Stats.UnladenRange),
					Inline: true,
				},
				{
					Name:   "Laden Jump Range",
					Value:  fmt.Sprint(t.Stats.LadenRange),
					Inline: true,
				},
				{
					Name:   "Total Unladen Jump Range",
					Value:  fmt.Sprint(t.Stats.AdditionalProperties["unladenFastestRange"]),
					Inline: true,
				},
				{
					Name:   "Total Laden Jump Rane",
					Value:  fmt.Sprint(t.Stats.AdditionalProperties["ladenFastestRange"]),
					Inline: true,
				},
				{
					Name:   "Mass lock Factor",
					Value:  fmt.Sprint(t.Stats.Masslock),
					Inline: true,
				},
				{
					Name:   "Shield",
					Value:  fmt.Sprint(t.Stats.Shield),
					Inline: true,
				},
				{
					Name:   "Integrity",
					Value:  fmt.Sprint(t.Stats.Armour),
					Inline: true,
				},
				{
					Name:   "Hardness",
					Value:  fmt.Sprint(t.Stats.AdditionalProperties["hardness"]),
					Inline: true,
				},
				{
					Name:   "DPS",
					Value:  fmt.Sprint(t.Stats.TotalDps),
					Inline: true,
				},
				{
					Name:   "EPS",
					Value:  fmt.Sprint(t.Stats.TotalEps),
					Inline: true,
				},
				{
					Name:   "Crew",
					Value:  fmt.Sprint(t.Stats.AdditionalProperties["crew"]),
					Inline: true,
				},
				{
					Name:   "Cargo",
					Value:  fmt.Sprint(t.Stats.CargoCapacity),
					Inline: true,
				},
				{
					Name:   "Passengers",
					Value:  fmt.Sprint(t.Stats.AdditionalProperties["passengerCapacity"]),
					Inline: true,
				},
				{
					Name:   "Fuel",
					Value:  fmt.Sprint(t.Stats.FuelCapacity),
					Inline: true,
				},
				{
					Name:   "Hull Mass",
					Value:  fmt.Sprint(t.Stats.HullMass),
					Inline: true,
				},
				{
					Name:   "Unladen Mass",
					Value:  fmt.Sprint(t.Stats.UnladenMass),
					Inline: true,
				},
				{
					Name:   "Laden Mass",
					Value:  fmt.Sprint(t.Stats.LadenMass),
					Inline: true,
				},

				{
					Name: "Core Internal",
					Value: "" +
						fmt.Sprintf("**%v**\n", t.Components.Standard.Bulkheads) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.PowerPlant) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.Thrusters) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.FrameShiftDrive) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.LifeSupport) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.PowerDistributor) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.Sensors) +
						fmt.Sprintf("**%v**\n", t.Components.Standard.FuelTank),
					Inline: true,
				},
				{
					Name:   "Optional Internal",
					Value:  optional,
					Inline: true,
				},
				{
					Name:   "Hardpoints",
					Value:  hard,
					Inline: true,
				},
				{
					Name:   "Utility",
					Value:  utils,
					Inline: true,
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("send Embed: %w", err)
	}

	return nil
}
