package inara

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
)

const URL = "https://inara.cz/inapi/v1/"

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

	r, err := api.Client.R().SetBody(data).Post(URL)
	if err != nil {
		return Profile{}, errors.Wrap(err, "inara getProfile")
	}
	var p Response
	err = json.Unmarshal(r.Body(), &p)
	if err != nil {
		return Profile{}, errors.Wrapf(err, "unmarshal inara profile: %v", string(r.Body()))
	}
	if p.Events[0].EventData.UserName == "" {
		return Profile{}, errors.New(fmt.Sprint("invalid Profile:", commander))
	}
	return p.Events[0].EventData, nil
}

func (api *API) GetProfiles(commanders []string) ([]Profile, error) {
	data := Data{
		Header: api.Header,
		Events: []Event{},
	}

	for i, commander := range commanders {
		data.Events = append(
			data.Events, Event{
				EventName:      "getCommanderProfile",
				EventTimestamp: time.Now().Format("2006-01-02T15:04:05Z"),
				EventCustomID:  i,
				EventData: struct {
					SearchName string `json:"searchName"`
				}{SearchName: commander},
			},
		)
	}

	r, err := api.Client.R().SetBody(data).Post(URL)
	if err != nil {
		return nil, errors.Wrap(err, "inara getProfile")
	}
	var p Response
	err = json.Unmarshal(r.Body(), &p)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal inara profile: %v", string(r.Body()))
	}
	for i, event := range p.Events {
		if event.EventData.UserName == "" {
			return nil, errors.New(fmt.Sprint("invalid Profile:", commanders[i]))
		}
	}
	var profiles []Profile
	From(p.Events).Select(func(i interface{}) interface{} { return i.(ResponseEvent).EventData }).ToSlice(&profiles)
	return profiles, nil
}
