package main

import "fmt"

func ConvertShipToURL(name string) string {
	buffer := ""
	buffer, ok := CoriolisShips[name]
	if !ok {
		buffer = EDSYShips[name]
	}

	return fmt.Sprintf("https://edassets.org/static/img/ship-schematics/qohen-leth/%v.png", buffer)
}

var CoriolisShips = map[string]string{
	"Adder":                "adder",
	"Alliance Challenger":  "alliance-challenger",
	"Alliance Chieftain":   "alliance-chieftain",
	"Alliance Crusader":    "alliance-crusader",
	"Anaconda":             "anaconda",
	"Asp Explorer":         "asp-explorer",
	"Asp Scout":            "asp-scout",
	"Beluga Liner":         "beluga-liner",
	"Cobra Mk III":         "cobra-mk-iii",
	"Cobra Mk IV":          "cobra-mk-iv",
	"Diamondback Explorer": "diamondback-explorer",
	"Diamondback Scout":    "diamondback-scout",
	"Dolphin":              "dolphin",
	"Eagle":                "eagle-mk-ii",
	"Federal Assault Ship": "federal-assault-ship",
	"Federal Corvette":     "federal-corvette",
	"Federal Dropship":     "federal-dropship",
	"Federal Gunship":      "federal-gunship",
	"Fer-de-Lance":         "fer-de-lance",
	"Hauler":               "hauler",
	"Imperial Clipper":     "imperial-clipper",
	"Imperial Courier":     "imperial-courier",
	"Imperial Cutter":      "imperial-cutter",
	"Imperial Eagle":       "imperial-eagle",
	"Keelback":             "keelback",
	"Krait Mk II":          "krait-mk-ii",
	"Krait Phantom":        "krait-phantom",
	"Mamba":                "mamba",
	"Orca":                 "orca",
	"Python":               "python",
	"Sidewinder":           "sidewinder",
	"Type-10 Defender":     "type-10-defender",
	"Type-6 Transporter":   "type-6-transporter",
	"Type-7 Transporter":   "type-7",
	"Type-9 Heavy":         "type-9-heavy",
	"Viper":                "viper-mk-iii",
	"Viper Mk IV":          "viper-mk-iv",
	"Vulture":              "vulture",
}

var EDSYShips = map[string]string{
	"Adder":                    "adder",
	"Anaconda":                 "anaconda",
	"Asp":                      "asp-explorer",
	"Asp_Scout":                "asp-scout",
	"BelugaLiner":              "beluga-liner",
	"CobraMkIII":               "cobra-mk-iii",
	"CobraMkIV":                "cobra-mk-iv",
	"Cutter":                   "imperial-cutter",
	"DiamondBackXL":            "diamondback-explorer",
	"DiamondBack":              "diamondback-scout",
	"Dolphin":                  "dolphin",
	"Eagle":                    "eagle-mk-ii",
	"Empire_Courier":           "imperial-courier",
	"Empire_Eagle":             "imperial-eagle",
	"Empire_Trader":            "imperial-clipper",
	"Federation_Corvette":      "federal-corvette",
	"Federation_Dropship":      "federal-dropship",
	"Federation_Dropship_MkII": "federal-assault-ship",
	"Federation_Gunship":       "federal-gunship",
	"FerDeLance":               "fer-de-lance",
	"Hauler":                   "hauler",
	"Independant_Trader":       "keelback",
	"Krait_MkII":               "krait-mk-ii",
	"Mamba":                    "mamba",
	"Krait_Light":              "krait-phantom",
	"Orca":                     "orca",
	"Python":                   "python",
	"SideWinder":               "sidewinder",
	"Type6":                    "type-6-transporter",
	"Type7":                    "type-7",
	"Type9":                    "type-9-heavy",
	"Type9_Military":           "type-10-defender",
	"TypeX":                    "alliance-chieftain",
	"TypeX_2":                  "alliance-crusader",
	"TypeX_3":                  "alliance-challenger",
	"Viper":                    "cobra-mk-iii",
	"Viper_MkIV":               "viper-mk-iv",
	"Vulture":                  "vulture",
}
