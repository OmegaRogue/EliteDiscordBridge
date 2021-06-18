package main

type InaraRequestData interface {
}

type InaraData struct {
	Header InaraHeader  `json:"header"`
	Events []InaraEvent `json:"events"`
}
type InaraHeader struct {
	AppName          string `json:"appName"`
	AppVersion       string `json:"appVersion"`
	IsBeingDeveloped bool   `json:"isBeingDeveloped"`
	APIkey           string `json:"APIkey"`
}
type InaraEvent struct {
	EventName      string           `json:"eventName"`
	EventTimestamp string           `json:"eventTimestamp"`
	EventCustomID  int              `json:"eventCustomID"`
	EventData      InaraRequestData `json:"eventData"`
}
