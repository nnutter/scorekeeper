package scorebook

import (
	"fmt"
	"strings"
	"time"
)

type Half string

const (
	Top          Half = "top"
	Bottom       Half = "bottom"
	BattingSlots      = 9
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
	Meta          GameMeta     `json:"meta"`
	Context       GameContext  `json:"context"`
	TopPitcher    string       `json:"topPitcher,omitempty"`
	BottomPitcher string       `json:"bottomPitcher,omitempty"`
	AwayOrder     []string     `json:"awayOrder,omitempty"`
	HomeOrder     []string     `json:"homeOrder,omitempty"`
	AwaySpot      int          `json:"awaySpot,omitempty"`
	HomeSpot      int          `json:"homeSpot,omitempty"`
	Entries       []EventEntry `json:"entries"`
}

func NewBook() Book {
	return Book{
		Context: GameContext{
			Inning: 1,
			Half:   Top,
		},
		AwayOrder: makeBattingOrder(),
		HomeOrder: makeBattingOrder(),
		Entries:   []EventEntry{},
	}
}

func (b *Book) AdvanceHalf() {
	b.SyncPitcherMemory()
	if b.Context.Half == Top {
		b.Context.Half = Bottom
		b.Context.Pitcher = b.rememberedPitcher(Bottom)
		return
	}
	b.Context.Half = Top
	b.Context.Inning++
	b.Context.Pitcher = b.rememberedPitcher(Top)
}

func (b *Book) RetreatHalf() {
	b.SyncPitcherMemory()
	if b.Context.Half == Bottom {
		b.Context.Half = Top
		b.Context.Pitcher = b.rememberedPitcher(Top)
		return
	}
	if b.Context.Inning > 1 {
		b.Context.Inning--
		b.Context.Half = Bottom
		b.Context.Pitcher = b.rememberedPitcher(Bottom)
	}
}

func (b *Book) SyncPitcherMemory() {
	switch b.Context.Half {
	case Bottom:
		b.BottomPitcher = b.Context.Pitcher
	default:
		b.TopPitcher = b.Context.Pitcher
	}
}

func (b *Book) HydratePitcherMemory() {
	b.TopPitcher = ""
	b.BottomPitcher = ""
	for _, entry := range b.Entries {
		switch entry.Half {
		case Bottom:
			b.BottomPitcher = entry.Pitcher
		default:
			b.TopPitcher = entry.Pitcher
		}
	}
	b.SyncPitcherMemory()
}

func (b *Book) HydrateBattingMemory() {
	b.AwayOrder = makeBattingOrder()
	b.HomeOrder = makeBattingOrder()
	b.AwaySpot = 0
	b.HomeSpot = 0

	for _, entry := range b.Entries {
		b.RecordPlateAppearance(entry)
	}
}

func (b *Book) HydrateMemory() {
	b.HydratePitcherMemory()
	b.HydrateBattingMemory()
}

func (b Book) rememberedPitcher(half Half) string {
	if half == Bottom {
		return b.BottomPitcher
	}
	return b.TopPitcher
}

func (b Book) RememberedBatter() string {
	order, spot := b.battingMemory(b.Context.Half)
	if spot < 0 || spot >= BattingSlots || spot >= len(order) {
		return ""
	}
	return order[spot]
}

func (b *Book) RecordPlateAppearance(entry EventEntry) {
	if entry.Mode != ModePlay {
		return
	}

	batter := strings.TrimSpace(entry.Batter)
	if batter == "" {
		return
	}

	order, spot := b.battingMemoryRef(entry.Half)
	*order = normalizeBattingOrder(*order)
	currentSpot := normalizeSpot(*spot)
	index := indexOfBatter(*order, batter)
	if index >= 0 {
		*spot = normalizeSpot(index + 1)
		return
	}

	(*order)[currentSpot] = batter
	*spot = normalizeSpot(currentSpot + 1)
}

func (b Book) battingMemory(half Half) ([]string, int) {
	if half == Bottom {
		return b.HomeOrder, b.HomeSpot
	}
	return b.AwayOrder, b.AwaySpot
}

func (b *Book) battingMemoryRef(half Half) (*[]string, *int) {
	if half == Bottom {
		return &b.HomeOrder, &b.HomeSpot
	}
	return &b.AwayOrder, &b.AwaySpot
}

func indexOfBatter(order []string, batter string) int {
	for i, candidate := range order {
		if candidate == batter {
			return i
		}
	}
	return -1
}

func makeBattingOrder() []string {
	return make([]string, BattingSlots)
}

func normalizeBattingOrder(order []string) []string {
	if len(order) == BattingSlots {
		return order
	}
	normalized := makeBattingOrder()
	copy(normalized, order)
	return normalized
}

func normalizeSpot(spot int) int {
	if BattingSlots == 0 {
		return 0
	}
	spot %= BattingSlots
	if spot < 0 {
		spot += BattingSlots
	}
	return spot
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
