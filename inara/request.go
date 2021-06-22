package inara

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type RequestData interface {
}

type Data struct {
	Header Header  `json:"header"`
	Events []Event `json:"events"`
}
type Header struct {
	AppName          string `json:"appName"`
	AppVersion       string `json:"appVersion"`
	IsBeingDeveloped bool   `json:"isBeingDeveloped"`
	APIkey           string `json:"APIkey"`
}
type Event struct {
	EventName      string      `json:"eventName"`
	EventTimestamp string      `json:"eventTimestamp"`
	EventCustomID  int         `json:"eventCustomID"`
	EventData      RequestData `json:"eventData"`
}

type API struct {
	Header Header
	Client *resty.Client
}

func NewAPI(header Header) *API {
	client := resty.New()
	api := API{
		Header: header,
		Client: client,
	}
	return &api
}

func (api *API) GetProfile(commander string) (Profile, error) {
	data := Data{
		Header: api.Header,
		Events: []Event{
			{
				EventName:      "getCommanderProfile",
				EventTimestamp: time.Now().Format("2006-01-02T15:04:05Z"),
				EventCustomID:  1234,
				EventData: struct {
					SearchName string `json:"searchName"`
				}{SearchName: commander},
			},
		},
	}

	r, err := api.Client.R().SetBody(data).Post("https://inara.cz/inapi/v1/")
	if err != nil {
		return Profile{}, fmt.Errorf("inara getProfile: %w", err)
	}
	var p Response
	err = json.Unmarshal(r.Body(), &p)
	if err != nil {
		return Profile{}, fmt.Errorf("unmarshal inara profile: %w\n%v", err, string(r.Body()))
	}
	if p.Events[0].EventData.UserName == "" {
		return Profile{}, fmt.Errorf("invalid Profile: %v", commander)
	}
	return p.Events[0].EventData, nil
}
