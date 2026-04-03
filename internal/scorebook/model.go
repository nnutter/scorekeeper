package scorebook

import (
	"fmt"
	"strings"
	"time"
)

type Half string

const (
	Top    Half = "top"
	Bottom Half = "bottom"
)

type EntryMode string

const (
	ModePlay EntryMode = "play"
	ModeRun  EntryMode = "run"
)

type GameMeta struct {
	AwayTeam string `json:"awayTeam"`
	HomeTeam string `json:"homeTeam"`
	GameDate string `json:"gameDate"`
}

type GameContext struct {
	Inning  int    `json:"inning"`
	Half    Half   `json:"half"`
	Pitcher string `json:"pitcher"`
}

type EventEntry struct {
	ID          string    `json:"id"`
	Mode        EntryMode `json:"mode"`
	Inning      int       `json:"inning"`
	Half        Half      `json:"half"`
	Pitcher     string    `json:"pitcher"`
	Batter      string    `json:"batter"`
	Pitches     string    `json:"pitches,omitempty"`
	BatterEvent string    `json:"batterEvent,omitempty"`
	Advances    string    `json:"advances,omitempty"`
	RunnerEvent string    `json:"runnerEvent,omitempty"`
	Note        string    `json:"note,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type EventDraft struct {
	EditingID   string
	Batter      string
	Pitches     string
	BatterEvent string
	Advances    string
	RunnerEvent string
	Note        string
}

type Book struct {
	Meta    GameMeta     `json:"meta"`
	Context GameContext  `json:"context"`
	Entries []EventEntry `json:"entries"`
}

func NewBook() Book {
	return Book{
		Context: GameContext{
			Inning: 1,
			Half:   Top,
		},
		Entries: []EventEntry{},
	}
}

func (b *Book) AdvanceHalf() {
	if b.Context.Half == Top {
		b.Context.Half = Bottom
		return
	}
	b.Context.Half = Top
	b.Context.Inning++
}

func (d EventDraft) ToEntry(ctx GameContext) EventEntry {
	mode := ModePlay
	if strings.TrimSpace(d.RunnerEvent) != "" && strings.TrimSpace(d.BatterEvent) == "" {
		mode = ModeRun
	}

	return EventEntry{
		ID:          nextID(),
		Mode:        mode,
		Inning:      ctx.Inning,
		Half:        ctx.Half,
		Pitcher:     strings.TrimSpace(ctx.Pitcher),
		Batter:      strings.TrimSpace(d.Batter),
		Pitches:     strings.TrimSpace(d.Pitches),
		BatterEvent: strings.TrimSpace(d.BatterEvent),
		Advances:    strings.TrimSpace(d.Advances),
		RunnerEvent: strings.TrimSpace(d.RunnerEvent),
		Note:        strings.TrimSpace(d.Note),
		CreatedAt:   time.Now(),
	}
}

func (d *EventDraft) LoadFromEntry(e EventEntry) {
	d.EditingID = e.ID
	d.Batter = e.Batter
	d.Pitches = e.Pitches
	d.BatterEvent = e.BatterEvent
	d.Advances = e.Advances
	d.RunnerEvent = e.RunnerEvent
	d.Note = e.Note
}

func (d *EventDraft) Reset() {
	d.EditingID = ""
	d.Batter = ""
	d.Pitches = ""
	d.BatterEvent = ""
	d.Advances = ""
	d.RunnerEvent = ""
	d.Note = ""
}

func (d *EventDraft) PrepareForNextRunnerEvent() {
	d.EditingID = ""
	d.BatterEvent = ""
	d.Advances = ""
	d.RunnerEvent = ""
	d.Note = ""
}

func (d *EventDraft) PrepareForNextPlateAppearance() {
	d.Reset()
}

func (e EventEntry) Summary() string {
	if e.Mode == ModeRun {
		return fmt.Sprintf("%s | batter %s | %s", strings.Title(string(e.Half)), e.Batter, e.RunnerEvent)
	}
	summary := fmt.Sprintf("%s | batter %s | %s", strings.Title(string(e.Half)), e.Batter, e.BatterEvent)
	if e.Advances != "" {
		summary += " | " + e.Advances
	}
	return summary
}

func nextID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
