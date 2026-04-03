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
