package scorebook

var PitchTokenRows = [][]string{
	[]string{"B", "C", "S", "K", "F", "L", "T"},
	[]string{"1", "2", "3", ">", "H", "I", "A"},
	{"+", "*", ".", "M", "N", "O", "P"},
	{"Q", "R", "U", "V", "X", "Y"},
}

var BatterTokenRows = [][]string{
	[]string{"K", "S", "D", "T", "HR", "DGR"},
	[]string{"W", "IW", "HP", "E", "FC", "FLE"},
	[]string{"(", ")", "1", "2", "3", "4", "5"},
	{"6", "7", "8", "9", ".", ";", "/"},
	{"SF", "SH", "GDP", "LDP", "FO"},
	[]string{"G", "L", "P", "F"},
}

var AdvanceTokenRows = [][]string{
	{"B-1", "B-2", "B-3", "B-H"},
	[]string{"1-2", "1-3", "1-H", "1X2", "1X3", "1XH"},
	[]string{"2-3", "2-H", "2X3", "2XH"},
	[]string{"3-H", "3XH"},
	{";", "(", ")", "E", "TH", "UR", "RBI", "NR"},
}

var RunnerTokenRows = [][]string{
	{"SB2", "SB3", "SBH", "WP"},
	[]string{"CS2", "CS3", "CSH", "BK"},
	[]string{"PO1", "PO2", "PO3", "DI"},
	[]string{"POCS2", "POCS3", "POCSH", "OA"},
	{"PB", ".", ";", "(", ")"},
}

const pitchChars = "+*.123>ABCFHIKLMNOPQRSTUVWXYZ"

func ValidPitchString(s string) bool {
	for _, r := range s {
		if r == ' ' {
			continue
		}
		if !containsRune(pitchChars, r) {
			return false
		}
	}
	return true
}

func containsRune(set string, target rune) bool {
	for _, r := range set {
		if r == target {
			return true
		}
	}
	return false
}
