package main

type Allegiance int

const (
	Independent Allegiance = iota
	Alliance
	Empire
	Federation
	Pirate
	PilotsFederation
	Thargoids
	Guardians
)

func ParseAllegiance(s string) Allegiance {
	allegiance, ok := allegianceMap[s]
	if !ok {
		return Independent
	}
	return allegiance
}

var allegianceMap = map[string]Allegiance{
	"Independent":      Independent,
	"Alliance":         Alliance,
	"Empire":           Empire,
	"Federation":       Federation,
	"Pirate":           Pirate,
	"PilotsFederation": PilotsFederation,
	"Thargoids":        Thargoids,
	"Guardians":        Guardians,
}
