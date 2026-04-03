package scorebook

import "testing"

func TestValidatePlayDraft(t *testing.T) {
	meta := GameMeta{AwayTeam: "Away", HomeTeam: "Home", GameDate: "2026-04-01"}
	ctx := GameContext{Inning: 1, Half: Top, Pitcher: "45S"}
	draft := EventDraft{Batter: "12J", Pitches: "CBX", BatterEvent: "S7"}

	issues := Validate(meta, ctx, draft)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestValidateRequiresOneEventType(t *testing.T) {
	meta := GameMeta{AwayTeam: "Away", HomeTeam: "Home", GameDate: "2026-04-01"}
	ctx := GameContext{Inning: 1, Half: Top, Pitcher: "45S"}
	draft := EventDraft{Batter: "12J"}

	issues := Validate(meta, ctx, draft)
	if len(issues) == 0 {
		t.Fatal("expected validation issue for missing event")
	}
	if issues[0] != "Enter a batter event or a base-running event." {
		t.Fatalf("unexpected issue: %v", issues)
	}
}

func TestValidateRejectsBothEventTypes(t *testing.T) {
	meta := GameMeta{AwayTeam: "Away", HomeTeam: "Home", GameDate: "2026-04-01"}
	ctx := GameContext{Inning: 1, Half: Top, Pitcher: "45S"}
	draft := EventDraft{Batter: "12J", BatterEvent: "S7", RunnerEvent: "SB2"}

	issues := Validate(meta, ctx, draft)
	if len(issues) == 0 {
		t.Fatal("expected validation issue when both event types are set")
	}
	if issues[0] != "Enter either a batter event or a base-running event, not both." {
		t.Fatalf("unexpected issue: %v", issues)
	}
}

func TestValidateFlagsInvalidPitchCharacters(t *testing.T) {
	meta := GameMeta{AwayTeam: "Away", HomeTeam: "Home", GameDate: "2026-04-01"}
	ctx := GameContext{Inning: 1, Half: Top, Pitcher: "45S"}
	draft := EventDraft{Batter: "12J", Pitches: "CBx", BatterEvent: "S7"}

	issues := Validate(meta, ctx, draft)
	if len(issues) == 0 {
		t.Fatal("expected invalid pitch syntax issue")
	}
}
