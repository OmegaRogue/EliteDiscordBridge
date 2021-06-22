package elite

type SLEF struct {
	Header struct {
		AppName    string `json:"appName"`
		AppVersion int    `json:"appVersion"`
		AppURL     string `json:"appURL"`
	} `json:"header"`
	Data struct {
		Event         string  `json:"event"`
		Ship          string  `json:"Ship"`
		ShipName      string  `json:"ShipName"`
		ShipIdent     string  `json:"ShipIdent"`
		HullValue     int     `json:"HullValue"`
		ModulesValue  int     `json:"ModulesValue"`
		UnladenMass   float64 `json:"UnladenMass"`
		CargoCapacity int     `json:"CargoCapacity"`
		MaxJumpRange  float64 `json:"MaxJumpRange"`
		FuelCapacity  struct {
			Main    int     `json:"Main"`
			Reserve float64 `json:"Reserve"`
		} `json:"FuelCapacity"`
		Rebuy   int `json:"Rebuy"`
		Modules []struct {
			Slot        string `json:"Slot"`
			Item        string `json:"Item"`
			On          bool   `json:"On"`
			Priority    int    `json:"Priority"`
			Value       int    `json:"Value,omitempty"`
			Engineering struct {
				BlueprintName string  `json:"BlueprintName"`
				Level         int     `json:"Level"`
				Quality       float64 `json:"Quality"`
				Modifiers     []struct {
					Label         string  `json:"Label"`
					Value         float64 `json:"Value"`
					OriginalValue float64 `json:"OriginalValue"`
				} `json:"Modifiers"`
			} `json:"Engineering,omitempty"`
		} `json:"Modules"`
	} `json:"data"`
}
