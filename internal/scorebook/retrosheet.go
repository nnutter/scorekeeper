package scorebook

var PitchTokenRows = [][]string{
	{"B", "C", "S", "K", "F", "X"},
	{"1", "2", "3", "A", "H", "I"},
	{">", "L", "T"},
	{"+", "*", ".", "M", "N", "O"},
	{"Q", "R", "U", "V", "Y", "P"},
}

var BatterTokenRows = [][]string{
	{"K", "S", "D", "T", "HR", "DGR"},
	{"W", "IW", "HP", "E", "FC", "FLE"},
	{"(", ")", "1", "2", "3", "4"},
	{"5", "6", "7", "8", "9", "/"},
	{"SF", "SH", "GDP", "LDP", "FO"},
	{"G", "L", "P", "F"},
}

var AdvanceTokenRows = [][]string{
	{"B-1", "B-2", "B-3", "B-H"},
	{"1-2", "1-3", "1-H", "1X2", "1X3", "1XH"},
	{"2-3", "2-H", "2X3", "2XH"},
	{"3-H", "3XH"},
	{"(", ")", "E", "TH", "UR", "RBI"},
}

var RunnerTokenRows = [][]string{
	{"SB2", "SB3", "SBH", "WP", "PB"},
	{"CS2", "CS3", "CSH", "BK"},
	{"PO1", "PO2", "PO3", "DI"},
	{"POCS2", "POCS3", "POCSH", "OA"},
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
