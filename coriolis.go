package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"edDiscord/coriolis"
	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/cdproto"
	"github.com/chromedp/chromedp"
)

func ShipBuildCoriolis(ctx context.Context, content string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	buildURL, err := url.Parse(content)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	allocatorContext, cancel := chromedp.NewRemoteAllocator(ctx, "ws://172.27.208.1:9223")
	defer cancel()
	ctx, cancel = chromedp.NewContext(
		allocatorContext,
		chromedp.WithLogf(devToolHandler),
		chromedp.WithErrorf(devToolHandler),
		chromedp.WithDebugf(devToolHandler),
		chromedp.WithBrowserOption(
			chromedp.WithBrowserDebugf(devToolHandler),
			chromedp.WithBrowserDebugf(devToolHandler),
			chromedp.WithBrowserLogf(devToolHandler),
		),
	)
	defer cancel()

	res := ""
	err = chromedp.Run(
		ctx,
		chromedp.Navigate(buildURL.String()),
		chromedp.Evaluate(
			"console.stdlog = console.log.bind(console);console.logs = [];console.log = function(){console.logs.push(Array.from(arguments));console.stdlog.apply(console, arguments);}",
			nil,
		),

		chromedp.WaitReady("#outfit"),
		chromedp.Evaluate("", &res),

		// chromedp.EvaluateAsDevTools(
		// 	"document.querySelector('#build > button:nth-child(7)').click();document.querySelector('textarea.cb').textContent",
		// 	&res,
		// ),
	)
	if err != nil {
		return fmt.Errorf("browser: %w", err)
	}

	fmt.Println(res)

	time.Sleep(time.Minute * 5)

	var t coriolis.ShipLoadout
	t.UnmarshalJSON([]byte(res))
	s.ChannelMessageDelete(m.ChannelID, m.ID)

	utils := "​"
	for _, v := range t.Components.Utility {

		if v != nil {
			v2 := v.(map[string]interface{})
			utils += fmt.Sprintf("**%v%v %v**\n", v2["class"], v2["rating"], v2["group"])
		}

	}
	optional := "​"
	for _, v := range t.Components.Internal {

		if v != nil {
			v2 := v.(map[string]interface{})
			optional += fmt.Sprintf("**%v%v %v**\n", v2["class"], v2["rating"], v2["group"])
		}

	}
	hard := "​"
	for _, v := range t.Components.Hardpoints {
		if v != nil {
			v2 := v.(map[string]interface{})
			hard += fmt.Sprintf("**%v%v %v %v**\n", v2["class"], v2["rating"], v2["mount"], v2["group"])
		}

	}

	_, err = s.ChannelMessageSendEmbed(
		m.ChannelID, &discordgo.MessageEmbed{
			URL:         buildURL.String(),
			Type:        discordgo.EmbedTypeRich,
			Title:       t.Name + "​",
			Description: fmt.Sprintf("%v​", t.Ship),
			Timestamp:   "",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf(
					"https://edassets.org/static/img/ship-schematics/qohen-leth/%v.png",
					CoriolisShips[t.Ship],
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
					Name:   "Laden Jump Range",
					Value:  fmt.Sprint(t.Stats.LadenRange),
					Inline: true,
				},
				{
					Name:   "Total Laden Range",
					Value:  fmt.Sprint(t.Stats.AdditionalProperties["ladenFastestRange"]),
					Inline: true,
				},
				{
					Name:   "​",
					Value:  "​",
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
					Name:   "​",
					Value:  "​",
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

func devToolHandler(s string, is ...interface{}) {
	/*
	   Uncomment the following line to have a log of the events
	   log.Printf(s, is...)
	*/
	/*
	   We need this to be on a separate gorutine
	   otherwise we block the browser and we don't receive messages
	*/
	go func() {
		for _, elem := range is {
			var msg cdproto.Message
			// The CDP messages are sent as strings so we need to convert them back
			err := json.Unmarshal([]byte(fmt.Sprintf("%s", elem)), &msg)
			// possible source of empty msg!!!!!!!!!!!!!
			if err != nil {
				log.Println(err)
				log.Printf("Faulty element:\n%v\n", fmt.Sprintf("%s", elem))
			}

		}
	}()
}
