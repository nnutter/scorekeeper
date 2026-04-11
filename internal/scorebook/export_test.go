package scorebook

import (
	"strings"
	"testing"
	"time"
)

func TestExportTextIncludesNotes(t *testing.T) {
	book := Book{
		Meta:    GameMeta{AwayTeam: "Away Club", HomeTeam: "Home Club", GameDate: "2026-04-01"},
		Context: GameContext{Inning: 1, Half: Top, Pitcher: "45S"},
		Entries: []EventEntry{
			{
				ID:          "1",
				Mode:        ModePlay,
				Inning:      1,
				Half:        Top,
				Pitcher:     "45S",
				Batter:      "12J",
				Pitches:     "CBX",
				BatterEvent: "S7",
				Advances:    "1-3",
				Note:        "lined into left",
				CreatedAt:   time.Unix(1, 0),
			},
			{
				ID:          "combo",
				Mode:        ModePlay,
				Inning:      1,
				Half:        Top,
				Pitcher:     "45S",
				Batter:      "13K",
				BatterEvent: "K",
				RunnerEvent: "SB2",
				CreatedAt:   time.Unix(1, 500),
			},
			{
				ID:          "2",
				Mode:        ModeRun,
				Inning:      1,
				Half:        Top,
				Pitcher:     "45S",
				Batter:      "12J",
				RunnerEvent: "SB2",
				CreatedAt:   time.Unix(2, 0),
			},
		},
	}

	export := ExportText(book)
	checks := []string{
		"game,2026-04-01,away=Away Club,home=Home Club",
		"play,1,top,pitcher=45S,batter=12J,pitches=CBX,event=S7,adv=1-3",
		"play,1,top,pitcher=45S,batter=13K,event=K+SB2",
		"note,1,top,\"lined into left\"",
		"run,1,top,pitcher=45S,batter=12J,event=SB2",
	}
	for _, check := range checks {
		if !strings.Contains(export, check) {
			t.Fatalf("export missing %q:\n%s", check, export)
		}
	}
}

func TestMailtoLinkEscapesBody(t *testing.T) {
	book := NewBook()
	book.Meta = GameMeta{AwayTeam: "Away Club", HomeTeam: "Home Club", GameDate: "2026-04-01"}

	link := MailtoLink(book)
	if !strings.HasPrefix(link, "mailto:?") {
		t.Fatalf("unexpected mailto link: %s", link)
	}
	if !strings.Contains(link, "subject=") || !strings.Contains(link, "body=") {
		t.Fatalf("mailto link missing subject/body: %s", link)
	}
}
