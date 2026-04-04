package scorebook

var PitchTokenRows = [][]string{
	{"+", "*", ".", "1", "2", "3", ">"},
	{"A", "B", "C", "F", "H", "I", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "Y"},
}

var BatterTokenRows = [][]string{
	{"S", "D", "T", "HR", "DGR", "K", "W", "IW", "HP", "E", "FC", "FLE", "NP"},
	{"SF", "SH", "GDP", "LDP", "FO", "G", "L", "P", "F", "1", "2", "3", "4", "5"},
	{"6", "7", "8", "9", "(", ")", ".", ";", "/"},
}

var AdvanceTokenRows = [][]string{
	{"B-1", "B-2", "B-3", "B-H", "1-2", "1-3", "1-H", "2-2", "2-3", "2-H", "3-3", "3-H"},
	{"1X2", "1X3", "1XH", "2X3", "2XH", "3XH", ";", "(", ")", "E", "TH", "UR", "RBI", "NR"},
}

var RunnerTokenRows = [][]string{
	{"SB2", "SB3", "SBH", "CS2", "CS3", "CSH", "PO1", "PO2", "PO3", "POCS2", "POCS3", "POCSH"},
	{"WP", "PB", "BK", "DI", "OA", ".", ";", "(", ")"},
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
