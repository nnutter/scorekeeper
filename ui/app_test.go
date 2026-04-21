package ui

import (
	"testing"

	"github.com/nnutter/scorekeeper/internal/scorebook"
)

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

func TestPitchCountLabelUsesCurrentPitcherAndDraft(t *testing.T) {
	r := &Root{
		book: scorebook.Book{
			Context: scorebook.GameContext{Pitcher: "45S"},
			Entries: []scorebook.EventEntry{
				{ID: "top-1", Pitcher: "45S", Pitches: "CB+"},
				{ID: "top-2", Pitcher: "12K", Pitches: "CBF"},
				{ID: "top-3", Pitcher: "45S", Pitches: "FX"},
			},
		},
		draft: scorebook.EventDraft{EditingID: "top-3", Pitches: "BV"},
	}

	if got := r.pitchCountLabel(); got != "P: 3" {
		t.Fatalf("pitchCountLabel() = %q, want %q", got, "P: 3")
	}
}

func TestSortLeadRunnerFirstAdvances(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"already sorted", "3-H;2-3;1-2;B-1", "3-H;2-3;1-2;B-1"},
		{"lead runner first", "1-2;3-H;B-1;2-3", "3-H;2-3;1-2;B-1"},
		{"outs included", "BX2;1X3;3XH;2-H", "3XH;2-H;1X3;BX2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortLeadRunnerFirst(tt.input, advanceSortRank)
			if got != tt.want {
				t.Fatalf("sortLeadRunnerFirst(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestUpsertAdvanceToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		token string
		want  string
	}{
		{"append first advance", "", "1-3", "1-3"},
		{"replace same runner", "1-3", "1X2", "1X2"},
		{"replace within list", "3-H;1-3;B-1", "1X2", "3-H;1X2;B-1"},
		{"append new runner", "2-3;1-2", "B-1", "2-3;1-2;B-1"},
		{"unknown token appends", "1-2", "E", "1-2;E"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := upsertAdvanceToken(tt.input, tt.token)
			if got != tt.want {
				t.Fatalf("upsertAdvanceToken(%q, %q) = %q, want %q", tt.input, tt.token, got, tt.want)
			}
		})
	}
}

func TestSortLeadRunnerFirstRunnerEvents(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"already sorted", "SBH;SB3;SB2", "SBH;SB3;SB2"},
		{"lead runner first", "SB2;PO3;SBH;PO1;SB3", "PO3;SBH;SB3;SB2;PO1"},
		{"picked off caught stealing", "POCS2;CSH;CS3", "CSH;CS3;POCS2"},
		{"unknown events stay last", "OA;SB3;WP;SB2", "SB3;SB2;OA;WP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortLeadRunnerFirst(tt.input, runnerEventSortRank)
			if got != tt.want {
				t.Fatalf("sortLeadRunnerFirst(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHandleHalfChangePreservesEditDraft(t *testing.T) {
	r := &Root{
		draft: scorebook.EventDraft{
			EditingID:   "123",
			Batter:      "12J",
			Pitches:     "CB",
			BatterEvent: "S7",
		},
		hasEditBase: true,
		message:     "Editing event.",
		messageKind: "status",
		focused:     "pitches",
		mobileKeys:  "advances",
	}

	r.handleHalfChange()

	if r.draft.EditingID != "123" || r.draft.Batter != "12J" || r.draft.Pitches != "CB" || r.draft.BatterEvent != "S7" {
		t.Fatalf("editing draft changed during half change: %+v", r.draft)
	}
	if !r.hasEditBase {
		t.Fatal("edit base should be preserved during half change")
	}
	if r.focused != "pitches" || r.mobileKeys != "advances" {
		t.Fatalf("editing UI state changed during half change: focused=%q mobileKeys=%q", r.focused, r.mobileKeys)
	}
	if r.message != "" || r.messageKind != "" {
		t.Fatalf("message should be cleared during half change: %q %q", r.message, r.messageKind)
	}
}

func TestHandleHalfChangeResetsNewEntryDraft(t *testing.T) {
	r := &Root{
		book: scorebook.Book{
			AwayOrder: make([]string, scorebook.BattingSlots),
			HomeOrder: make([]string, scorebook.BattingSlots),
		},
		draft: scorebook.EventDraft{
			Batter:      "12J",
			Pitches:     "CB",
			BatterEvent: "S7",
		},
		hasEditBase: true,
		message:     "Event saved.",
		messageKind: "status",
		focused:     "pitches",
		mobileKeys:  "advances",
	}

	r.handleHalfChange()

	if r.draft != (scorebook.EventDraft{}) {
		t.Fatalf("new-entry draft should reset during half change: %+v", r.draft)
	}
	if r.hasEditBase {
		t.Fatal("edit base should be cleared for new-entry half change")
	}
	if r.focused != "" || r.mobileKeys != "pitches" {
		t.Fatalf("new-entry UI state not reset: focused=%q mobileKeys=%q", r.focused, r.mobileKeys)
	}
	if r.message != "" || r.messageKind != "" {
		t.Fatalf("message should be cleared during half change: %q %q", r.message, r.messageKind)
	}
}

func TestStepBatterUpdatesEditingDraftBatter(t *testing.T) {
	r := &Root{
		book: scorebook.Book{
			AwayOrder: []string{"A1", "A2", "A3", "", "", "", "", "", ""},
			Context:   scorebook.GameContext{Half: scorebook.Top},
		},
		draft:      scorebook.EventDraft{EditingID: "top-1", Batter: "A1"},
		editBatter: 1,
	}

	r.stepBatter(1)
	if r.editBatter != 2 {
		t.Fatalf("edit batter position after advance = %d, want 2", r.editBatter)
	}
	if r.draft.Batter != "A2" {
		t.Fatalf("editing batter after advance = %q, want A2", r.draft.Batter)
	}

	r.stepBatter(-1)
	if r.editBatter != 1 {
		t.Fatalf("edit batter position after retreat = %d, want 1", r.editBatter)
	}
	if r.draft.Batter != "A1" {
		t.Fatalf("editing batter after retreat = %q, want A1", r.draft.Batter)
	}
}

func TestStepBatterWrapsWhileEditing(t *testing.T) {
	r := &Root{
		book: scorebook.Book{
			HomeOrder: []string{"H1", "", "", "", "", "", "", "", "H9"},
			Context:   scorebook.GameContext{Half: scorebook.Bottom},
		},
		draft:      scorebook.EventDraft{EditingID: "bottom-1", Batter: "H1"},
		editBatter: 1,
	}

	r.stepBatter(-1)
	if r.editBatter != 9 {
		t.Fatalf("edit batter position after wrap = %d, want 9", r.editBatter)
	}
	if r.draft.Batter != "H9" {
		t.Fatalf("editing batter after wrap = %q, want H9", r.draft.Batter)
	}
}

func TestCurrentEntryBattingPositionUsesEditBatterWhileEditing(t *testing.T) {
	r := &Root{
		book:       scorebook.Book{Context: scorebook.GameContext{Half: scorebook.Top}, AwaySpot: 1},
		draft:      scorebook.EventDraft{EditingID: "top-1"},
		editBatter: 4,
	}

	if got := r.currentEntryBattingPosition(); got != 4 {
		t.Fatalf("current entry batting position = %d, want 4", got)
	}
}

func TestCurrentEntryBattingPositionFallsBackToBookPosition(t *testing.T) {
	r := &Root{
		book: scorebook.Book{
			Context:  scorebook.GameContext{Half: scorebook.Top},
			AwaySpot: 2,
		},
	}

	if got := r.currentEntryBattingPosition(); got != 3 {
		t.Fatalf("current entry batting position = %d, want 3", got)
	}
}

func TestSortedLogEntriesOrdersByInningHalfAndBattingPositionStably(t *testing.T) {
	book := scorebook.Book{
		Context: scorebook.GameContext{Inning: 2, Half: scorebook.Bottom},
		Entries: []scorebook.EventEntry{
			{ID: "top-1", Inning: 1, Half: scorebook.Top, Mode: scorebook.ModePlay, Batter: "A1", BattingPos: 1},
			{ID: "bottom-1", Inning: 1, Half: scorebook.Bottom, Mode: scorebook.ModePlay, Batter: "H1", BattingPos: 1},
			{ID: "top-2", Inning: 1, Half: scorebook.Top, Mode: scorebook.ModePlay, Batter: "A2", BattingPos: 2},
			{ID: "bottom-2", Inning: 1, Half: scorebook.Bottom, Mode: scorebook.ModePlay, Batter: "H2", BattingPos: 2},
			{ID: "top-3-a", Inning: 2, Half: scorebook.Top, Mode: scorebook.ModePlay, Batter: "A3", BattingPos: 3},
			{ID: "top-3-b", Inning: 2, Half: scorebook.Top, Mode: scorebook.ModeRun, Batter: "A2", RunnerEvent: "SB2", BattingPos: 3},
		},
	}

	sorted := sortedLogEntries(book)
	want := []string{"top-1", "top-2", "bottom-1", "bottom-2", "top-3-a", "top-3-b"}
	for i, id := range want {
		if sorted[i].ID != id {
			t.Fatalf("sorted[%d] = %q, want %q", i, sorted[i].ID, id)
		}
	}

	if book.Entries[0].ID != "top-1" {
		t.Fatalf("sortedLogEntries should not mutate input, first entry = %q", book.Entries[0].ID)
	}
}

func TestSortedLogEntriesReordersEditedEntryWhenBattingPositionCollides(t *testing.T) {
	book := scorebook.Book{
		Entries: []scorebook.EventEntry{
			{ID: "top-1", Inning: 1, Half: scorebook.Top, Mode: scorebook.ModePlay, Batter: "A1", BattingPos: 2},
			{ID: "top-2", Inning: 1, Half: scorebook.Top, Mode: scorebook.ModePlay, Batter: "A2"},
		},
	}

	sorted := sortedLogEntries(book)
	want := []string{"top-2", "top-1"}
	for i, id := range want {
		if sorted[i].ID != id {
			t.Fatalf("sorted[%d] = %q, want %q", i, sorted[i].ID, id)
		}
	}
}
