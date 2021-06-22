package elite

type Power int

const (
	Uncontrolled Power = iota
	AislingDuval
	ArchonDelaine
	ArissaLavignyDuval
	DentonPatreus
	EdmundMahon
	FeliciaWinters
	LiYongRui
	PranavAntal
	YuriGrom
	ZacharyHudson
	ZeminaTorval
)

func ParsePower(s string) Power {
	power, ok := powerMap[s]
	if !ok {
		return Uncontrolled
	}
	return power
}

var powerMap = map[string]Power{
	"Aisling Duval":        AislingDuval,
	"Archon Delaine":       ArchonDelaine,
	"Arissa Lavigny-Duval": ArissaLavignyDuval,
	"Denton Patreus":       DentonPatreus,
	"Edmund Mahon":         EdmundMahon,
	"Felicia Winters":      FeliciaWinters,
	"Li Yong-Rui":          LiYongRui,
	"Pranav Antal":         PranavAntal,
	"Yuri Grom":            YuriGrom,
	"Zachary Hudson":       ZacharyHudson,
	"Zemina Torval":        ZeminaTorval,
}
