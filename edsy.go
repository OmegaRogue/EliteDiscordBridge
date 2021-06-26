package main

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"edDiscord/edsy"
	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func ShipBuildEDSY(ctx context.Context, content string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	buildURL, err := url.Parse(content)
	if err != nil {
		return errors.Wrap(err, "parse url")
	}

	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	hash := strings.TrimPrefix(buildURL.Fragment, "/L=")

	var res edsy.Build
	err = chromedp.Run(
		ctx,
		chromedp.Navigate("https://edsy.org/api.html"),
		chromedp.Evaluate(
			fmt.Sprintf(
				"Build = window.edsy.Build;ship = Build.fromHash(\"%s\");ship.updateStats();ship.getHash();[\"component\",\"hardpoint\",\"internal\",\"military\",\"utility\"].map((a)=>{return ship.slots[a].map((c,i)=>{c.build=null;return c;});});ship.slots.ship.hull.build=null;ship.slots.ship.hatch.build=null;ship",
				hash,
			), &res,
		),
	)
	if err != nil {
		return errors.Wrap(err, "browser")
	}
	utils := "​"
	optional := "​"
	hard := "​"
	for _, v := range res.Slots.Utility {
		utils += fmt.Sprintln(v)
	}
	for _, v := range append(res.Slots.Internal, res.Slots.Military...) {
		optional += fmt.Sprintln(v)
	}
	for _, v := range res.Slots.Hardpoint {
		hard += fmt.Sprintln(v)
	}

	name := fmt.Sprintf("%s %s", res.Name, res.Nametag)
	if res.Name == "" {
		name = res.Slots.Ship.Hull.Module.Name
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return errors.Wrap(err, "delete url message")
	}

	_, err = s.ChannelMessageSendEmbed(
		m.ChannelID, &discordgo.MessageEmbed{
			URL:   buildURL.String(),
			Type:  discordgo.EmbedTypeRich,
			Title: name + "​",
			Description: fmt.Sprintf(
				"[%v](%v)​\n[Stations that sell this Build](%v)",
				res.Slots.Ship.Hull.Module.Name,
				res.Shipid,
				res.GetEDDBLink(),
			),
			Timestamp: "",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf(
					"https://edassets.org/static/img/ship-schematics/qohen-leth/%v.png",
					EDSYShips[res.Slots.Ship.Hull.Module.Fdname],
				),
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
					Value:  fmt.Sprintf("%.2f", res.Stats.Speed),
					Inline: true,
				},
				{
					Name:   "Boost",
					Value:  fmt.Sprintf("%.2f", res.Stats.Boost),
					Inline: true,
				},
				{
					Name:   "​",
					Value:  "​",
					Inline: true,
				},
				{
					Name:   "Laden Jump Range",
					Value:  fmt.Sprintf("%.2f", res.Stats.JumpLaden),
					Inline: true,
				},
				{
					Name:   "Total Laden Range",
					Value:  fmt.Sprintf("%.2f", res.Stats.RangeLaden),
					Inline: true,
				},
				{
					Name:   "​",
					Value:  "​",
					Inline: true,
				},
				{
					Name:   "Shield",
					Value:  fmt.Sprintf("%.2f", res.Stats.Shields),
					Inline: true,
				},
				{
					Name:   "Integrity",
					Value:  fmt.Sprintf("%.2f", res.Stats.Armour),
					Inline: true,
				},
				{
					Name:   "​",
					Value:  "​",
					Inline: true,
				},
				{
					Name:   "Cargo",
					Value:  fmt.Sprint(res.Stats.Cargocap),
					Inline: true,
				},
				{
					Name:   "Passengers",
					Value:  fmt.Sprint(res.Stats.Cabincap),
					Inline: true,
				},
				{
					Name:   "Fuel",
					Value:  fmt.Sprint(res.Stats.Fuelcap),
					Inline: true,
				},
				{
					Name:   "Core Internal",
					Value:  fmt.Sprintln(res.Slots.Component),
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
				{
					Name: "Cost",
					Value: fmt.Sprintf(
						"```diff\nCost:         %v\nRebuy Cost:   %v\nRestock Cost: %v\nRearm Cost:   %v\nVehicle Cost: %v\n```\n",
						strings.Replace(fmt.Sprintf("% 12d", res.Stats.Cost), " ", ".", -1),
						strings.Replace(fmt.Sprintf("% 12d", int(float32(res.Stats.Cost)*0.05)), " ", ".", -1),
						strings.Replace(fmt.Sprintf("% 12d", res.Stats.CostRestock), " ", ".", -1),
						strings.Replace(fmt.Sprintf("% 12d", res.Stats.CostRearm), " ", ".", -1),
						strings.Replace(fmt.Sprintf("% 12d", res.Stats.CostVehicle), " ", ".", -1),
					),
					Inline: false,
				},
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "send Embed")
	}

	return nil

}
