package eddb

import (
	"fmt"
	"strings"
)

func GetEDDBLink(shipID int, moduleID ...[]int) string {
	modules := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(moduleID)), ","), "[]")
	return fmt.Sprintf("https://eddb.io/station?s=%v&m=%v", shipID, modules)
}

func GetEDDBShip(shipID int) string {

	return fmt.Sprintf("https://eddb.io/station?s=%v", shipID)
}

func GetEDDBModule(moduleID int) string {

	return fmt.Sprintf("https://eddb.io/station?m=%v", moduleID)
}
