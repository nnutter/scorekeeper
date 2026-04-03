package scorebook

import "testing"

func TestValidPitchString(t *testing.T) {
	if !ValidPitchString("CBX") {
		t.Fatal("expected CBX to be valid")
	}
	if !ValidPitchString("B1>F") {
		t.Fatal("expected B1>F to be valid")
	}
	if ValidPitchString("CBx") {
		t.Fatal("expected lowercase x to be invalid")
	}
}
