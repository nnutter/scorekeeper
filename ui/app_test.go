package ui

import "testing"

func TestCountBallsStrikes(t *testing.T) {
	tests := []struct {
		name    string
		pitches string
		balls   int
		strikes int
	}{
		{"empty", "", 0, 0},
		{"single ball B", "B", 1, 0},
		{"single ball I", "I", 1, 0},
		{"single ball P", "P", 1, 0},
		{"single ball V", "V", 1, 0},
		{"single called strike A", "A", 0, 1},
		{"single called strike C", "C", 0, 1},
		{"single called strike K", "K", 0, 1},
		{"single called strike M", "M", 0, 1},
		{"single called strike Q", "Q", 0, 1},
		{"single called strike S", "S", 0, 1},
		{"foul when zero strikes", "F", 0, 1},
		{"foul when one strike", "SF", 0, 2},
		{"foul when two strikes", "SFF", 0, 2},
		{"foul L when one strike", "CL", 0, 2},
		{"foul L when two strikes", "CLR", 0, 2},
		{"other H ignored", "H", 0, 0},
		{"other O ignored", "O", 0, 0},
		{"other T ignored", "T", 0, 0},
		{"other U ignored", "U", 0, 0},
		{"other X ignored", "X", 0, 0},
		{"other Y ignored", "Y", 0, 0},
		{"other + ignored", "+", 0, 0},
		{"other * ignored", "*", 0, 0},
		{"other . ignored", ".", 0, 0},
		{"other 1 ignored", "1", 0, 0},
		{"other 2 ignored", "2", 0, 0},
		{"other 3 ignored", "3", 0, 0},
		{"other > ignored", ">", 0, 0},
		{"other N ignored", "N", 0, 0},
		{"mixed sequence", "BCFSK", 1, 4},
		{"realistic at-bat", "CBFBX", 2, 2},
		{"three balls two strikes", "BBCCF", 2, 2},
		{"all ball types", "BIPV", 4, 0},
		{"called strikes only", "ACKMQS", 0, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balls, strikes := countBallsStrikes(tt.pitches)
			if balls != tt.balls || strikes != tt.strikes {
				t.Fatalf("countBallsStrikes(%q) = %d-%d, want %d-%d", tt.pitches, balls, strikes, tt.balls, tt.strikes)
			}
		})
	}
}
