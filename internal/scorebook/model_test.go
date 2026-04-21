package scorebook

import "testing"

func TestAdvanceHalf(t *testing.T) {
	book := NewBook()

	book.AdvanceHalf()
	if book.Context.Inning != 1 || book.Context.Half != Bottom {
		t.Fatalf("first advance = inning %d half %s, want inning 1 bottom", book.Context.Inning, book.Context.Half)
	}

	book.AdvanceHalf()
	if book.Context.Inning != 2 || book.Context.Half != Top {
		t.Fatalf("second advance = inning %d half %s, want inning 2 top", book.Context.Inning, book.Context.Half)
	}
}

func TestRetreatHalf(t *testing.T) {
	book := NewBook()

	book.RetreatHalf()
	if book.Context.Inning != 1 || book.Context.Half != Top {
		t.Fatalf("retreat from start = inning %d half %s, want inning 1 top", book.Context.Inning, book.Context.Half)
	}

	book.AdvanceHalf()
	book.RetreatHalf()
	if book.Context.Inning != 1 || book.Context.Half != Top {
		t.Fatalf("retreat from inning 1 bottom = inning %d half %s, want inning 1 top", book.Context.Inning, book.Context.Half)
	}

	book.AdvanceHalf()
	book.AdvanceHalf()
	book.RetreatHalf()
	if book.Context.Inning != 1 || book.Context.Half != Bottom {
		t.Fatalf("retreat from inning 2 top = inning %d half %s, want inning 1 bottom", book.Context.Inning, book.Context.Half)
	}
}

func TestAdvanceBattingPositionWrapsWithinHalf(t *testing.T) {
	book := NewBook()

	book.AdvanceBattingPosition()
	if got := book.BattingPosition(); got != 2 {
		t.Fatalf("batting position after advance = %d, want 2", got)
	}

	book.AwaySpot = 8
	book.AdvanceBattingPosition()
	if got := book.BattingPosition(); got != 1 {
		t.Fatalf("batting position after wrap = %d, want 1", got)
	}

	book.Context.Half = Bottom
	if got := book.BattingPosition(); got != 1 {
		t.Fatalf("bottom batting position should be independent = %d, want 1", got)
	}
}

func TestRetreatBattingPositionWrapsWithinHalf(t *testing.T) {
	book := NewBook()

	book.RetreatBattingPosition()
	if got := book.BattingPosition(); got != 9 {
		t.Fatalf("batting position after retreat = %d, want 9", got)
	}

	book.Context.Half = Bottom
	book.HomeSpot = 3
	book.RetreatBattingPosition()
	if got := book.BattingPosition(); got != 3 {
		t.Fatalf("bottom batting position after retreat = %d, want 3", got)
	}

	book.Context.Half = Top
	if got := book.BattingPosition(); got != 9 {
		t.Fatalf("top batting position should be unchanged = %d, want 9", got)
	}
}

func TestHalfSwitchRestoresRememberedPitchers(t *testing.T) {
	book := NewBook()
	book.Context.Pitcher = "H1"

	book.AdvanceHalf()
	if book.Context.Pitcher != "" {
		t.Fatalf("advance to bottom pitcher = %q, want empty", book.Context.Pitcher)
	}

	book.Context.Pitcher = "A1"
	book.AdvanceHalf()
	if book.Context.Inning != 2 || book.Context.Half != Top || book.Context.Pitcher != "H1" {
		t.Fatalf("advance to inning 2 top = inning %d half %s pitcher %q, want inning 2 top pitcher H1", book.Context.Inning, book.Context.Half, book.Context.Pitcher)
	}

	book.RetreatHalf()
	if book.Context.Inning != 1 || book.Context.Half != Bottom || book.Context.Pitcher != "A1" {
		t.Fatalf("retreat to inning 1 bottom = inning %d half %s pitcher %q, want inning 1 bottom pitcher A1", book.Context.Inning, book.Context.Half, book.Context.Pitcher)
	}
}

func TestHydratePitcherMemory(t *testing.T) {
	book := Book{
		Context: GameContext{Inning: 3, Half: Bottom, Pitcher: "A2"},
		Entries: []EventEntry{
			{Inning: 1, Half: Top, Pitcher: "H1"},
			{Inning: 1, Half: Bottom, Pitcher: "A1"},
			{Inning: 2, Half: Top, Pitcher: "H2"},
		},
	}

	book.HydratePitcherMemory()

	if book.TopPitcher != "H2" {
		t.Fatalf("top pitcher memory = %q, want H2", book.TopPitcher)
	}
	if book.BottomPitcher != "A2" {
		t.Fatalf("bottom pitcher memory = %q, want A2", book.BottomPitcher)
	}
}

func TestCountPitchesIgnoresNonCountingCodes(t *testing.T) {
	tests := []struct {
		name    string
		pitches string
		want    int
	}{
		{"empty", "", 0},
		{"counts regular pitches", "CBFX", 4},
		{"ignores excluded codes", "+*.123>ANV", 0},
		{"mixed sequence", "C+BVNFX", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountPitches(tt.pitches); got != tt.want {
				t.Fatalf("CountPitches(%q) = %d, want %d", tt.pitches, got, tt.want)
			}
		})
	}
}

func TestPitchCountForPitcherUsesCurrentDraftAndSkipsEditedEntry(t *testing.T) {
	book := Book{
		Entries: []EventEntry{
			{ID: "top-1", Pitcher: "45S", Pitches: "CB+"},
			{ID: "top-2", Pitcher: "45S", Pitches: "ANV"},
			{ID: "top-3", Pitcher: "12K", Pitches: "CBF"},
			{ID: "top-4", Pitcher: "45S", Pitches: "FX"},
		},
	}

	if got := book.PitchCountForPitcher("45S", "", "BV"); got != 5 {
		t.Fatalf("pitch count for 45S = %d, want 5", got)
	}
	if got := book.PitchCountForPitcher("45S", "top-4", "BV"); got != 3 {
		t.Fatalf("pitch count for edited 45S entry = %d, want 3", got)
	}
	if got := book.PitchCountForPitcher("", "", "BV"); got != 0 {
		t.Fatalf("pitch count for blank pitcher = %d, want 0", got)
	}
}

func TestRecordPlateAppearanceLearnsBattingOrderByTeam(t *testing.T) {
	book := NewBook()

	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A1"})
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A2"})
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Bottom, Batter: "H1"})

	if got := book.RememberedBatter(); got != "" {
		t.Fatalf("top remembered batter = %q, want empty for unknown next slot", got)
	}

	book.Context.Half = Bottom
	if got := book.RememberedBatter(); got != "" {
		t.Fatalf("bottom remembered batter = %q, want empty for unknown next slot", got)
	}
}

func TestRecordPlateAppearanceIgnoresRunnerEvents(t *testing.T) {
	book := NewBook()
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A1"})
	book.RecordPlateAppearance(EventEntry{Mode: ModeRun, Half: Top, Batter: "A1", RunnerEvent: "SB2"})

	if got := book.RememberedBatter(); got != "" {
		t.Fatalf("remembered batter after runner event = %q, want empty for unknown next slot", got)
	}
}

func TestHydrateBattingMemoryRebuildsExpectedBatter(t *testing.T) {
	book := Book{
		Context: GameContext{Inning: 3, Half: Top},
		Entries: []EventEntry{
			{Mode: ModePlay, Inning: 1, Half: Top, Batter: "A1"},
			{Mode: ModePlay, Inning: 1, Half: Top, Batter: "A2"},
			{Mode: ModeRun, Inning: 1, Half: Top, Batter: "A2", RunnerEvent: "SB2"},
			{Mode: ModePlay, Inning: 1, Half: Bottom, Batter: "H1"},
			{Mode: ModePlay, Inning: 2, Half: Top, Batter: "A1"},
		},
	}

	book.HydrateBattingMemory()

	if got := book.RememberedBatter(); got != "A2" {
		t.Fatalf("top remembered batter after hydrate = %q, want A2", got)
	}

	book.Context.Half = Bottom
	if got := book.RememberedBatter(); got != "" {
		t.Fatalf("bottom remembered batter after hydrate = %q, want empty for unknown next slot", got)
	}
}

func TestRecordPlateAppearanceReturnsLearnedBatterAfterWrap(t *testing.T) {
	book := NewBook()
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A1"})
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A2"})
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A3"})

	if got := book.RememberedBatter(); got != "" {
		t.Fatalf("remembered batter before wrap = %q, want empty", got)
	}

	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A1"})

	if got := book.RememberedBatter(); got != "A2" {
		t.Fatalf("remembered batter after wrap = %q, want A2", got)
	}
}

func TestBattingPositionForEntryReturnsHistoricalSpot(t *testing.T) {
	book := Book{
		Context: GameContext{Inning: 3, Half: Bottom},
		Entries: []EventEntry{
			{ID: "top-1", Mode: ModePlay, Half: Top, Batter: "A1"},
			{ID: "top-2", Mode: ModePlay, Half: Top, Batter: "A2"},
			{ID: "bottom-1", Mode: ModePlay, Half: Bottom, Batter: "H1"},
			{ID: "top-run", Mode: ModeRun, Half: Top, Batter: "A2", RunnerEvent: "SB2"},
			{ID: "top-3", Mode: ModePlay, Half: Top, Batter: "A3"},
		},
	}

	if got := book.BattingPositionForEntry("top-2"); got != 2 {
		t.Fatalf("top-2 batting position = %d, want 2", got)
	}
	if got := book.BattingPositionForEntry("top-run"); got != 3 {
		t.Fatalf("top-run batting position = %d, want 3", got)
	}
	if got := book.BattingPositionForEntry("bottom-1"); got != 1 {
		t.Fatalf("bottom-1 batting position = %d, want 1", got)
	}
	if got := book.BattingPositionForEntry("missing"); got != 1 {
		t.Fatalf("missing batting position = %d, want current bottom spot 1", got)
	}
}

func TestBattingPositionForEntryPrefersStoredPosition(t *testing.T) {
	book := Book{
		Entries: []EventEntry{
			{ID: "top-1", Mode: ModePlay, Half: Top, Batter: "A1", BattingPos: 4},
		},
	}

	if got := book.BattingPositionForEntry("top-1"); got != 4 {
		t.Fatalf("stored batting position = %d, want 4", got)
	}
}

func TestRecordPlateAppearanceUsesStoredPosition(t *testing.T) {
	book := NewBook()
	book.RecordPlateAppearance(EventEntry{Mode: ModePlay, Half: Top, Batter: "A4", BattingPos: 4})

	if got := book.AwayOrder[3]; got != "A4" {
		t.Fatalf("stored batter slotted at 4 = %q, want A4", got)
	}
	if got := book.BattingPosition(); got != 5 {
		t.Fatalf("next batting position after stored slot = %d, want 5", got)
	}
}

func TestEventDraftToEntryTrimsValues(t *testing.T) {
	draft := EventDraft{
		Batter:      " 12J ",
		Pitches:     " CBX ",
		BatterEvent: " S7 ",
		Advances:    " 1-3 ",
		Note:        " loud contact ",
	}

	entry := draft.ToEntry(GameContext{Inning: 3, Half: Bottom, Pitcher: " 45S "})

	if entry.Inning != 3 || entry.Half != Bottom {
		t.Fatalf("unexpected context copied: %+v", entry)
	}
	if entry.Pitcher != "45S" || entry.Batter != "12J" || entry.Pitches != "CBX" || entry.BatterEvent != "S7" || entry.Advances != "1-3" || entry.Note != "loud contact" {
		t.Fatalf("unexpected trimmed entry: %+v", entry)
	}
}

func TestEventDraftToEntryDetectsRunnerMode(t *testing.T) {
	draft := EventDraft{
		Batter:      "12J",
		Pitches:     "CBX",
		RunnerEvent: "SB2",
	}

	entry := draft.ToEntry(GameContext{Inning: 1, Half: Top, Pitcher: "45S"})

	if entry.Mode != ModeRun {
		t.Fatalf("expected runner mode, got %s", entry.Mode)
	}
}

func TestPrepareForNextRunnerEventPreservesBatterAndPitches(t *testing.T) {
	draft := EventDraft{
		Batter:      "12J",
		Pitches:     "CBX",
		RunnerEvent: "SB2",
		Note:        "jumped early",
	}

	draft.PrepareForNextRunnerEvent()

	if draft.Batter != "12J" || draft.Pitches != "CBX" {
		t.Fatalf("expected batter and pitches preserved: %+v", draft)
	}
	if draft.RunnerEvent != "" || draft.Note != "" || draft.BatterEvent != "" || draft.Advances != "" {
		t.Fatalf("expected runner-only fields cleared: %+v", draft)
	}
}

func TestPrepareForNextPlateAppearanceClearsEntryFields(t *testing.T) {
	draft := EventDraft{
		Batter:      "12J",
		Pitches:     "CBX",
		BatterEvent: "S7",
		Advances:    "1-3",
		RunnerEvent: "",
		Note:        "lined out",
	}

	draft.PrepareForNextPlateAppearance()

	if draft.Batter != "" || draft.Pitches != "" || draft.BatterEvent != "" || draft.Advances != "" || draft.Note != "" {
		t.Fatalf("expected plate appearance fields cleared: %+v", draft)
	}
}

func TestEventEntryEventTextCombinesBatterAndRunnerEvents(t *testing.T) {
	entry := EventEntry{BatterEvent: "K", RunnerEvent: "SB2"}

	if got := entry.EventText(); got != "K+SB2" {
		t.Fatalf("event text = %q, want K+SB2", got)
	}
}

func TestEventEntryLogEventTextMatchesLogFormatting(t *testing.T) {
	entry := EventEntry{BatterEvent: "S7", RunnerEvent: "SB2", Advances: "1-3"}

	if got := entry.LogEventText(); got != "S7+SB2.1-3" {
		t.Fatalf("log event text = %q, want %q", got, "S7+SB2.1-3")
	}
}
