package kjspkg

import (
	"strconv"
)

var VersionsInOrder = []string{"1.12", "1.16", "1.18", "1.19", "1.20", "1.21"}
var Versions = map[string]int{
	"1.12": 2,
	"1.16": 6,
	"1.18": 8,
	"1.19": 9,
	"1.20": 10,
	"1.21": 11,
}

func GetVersionString(id int) string {
	for ver, thisid := range Versions {
		if thisid == id {
			return ver
		}
	}
	return "??? " + strconv.Itoa(id)
}
