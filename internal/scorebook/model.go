package scorebook

import (
	"fmt"
	"strings"
	"time"
)

type Half string

const (
	Top          Half = "▲"
	Bottom       Half = "▼"
	DefaultBattingSlots = 9
)

type EntryMode string

const (
	ModePlay EntryMode = "play"
	ModeRun  EntryMode = "run"
)

type GameMeta struct {
	AwayTeam  string `json:"awayTeam"`
	AwaySlots int    `json:"awaySlots,omitempty"`
	HomeTeam  string `json:"homeTeam"`
	HomeSlots int    `json:"homeSlots,omitempty"`
	GameDate  string `json:"gameDate"`
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
	BattingPos  int       `json:"battingPos,omitempty"`
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
		Meta: normalizedMeta(GameMeta{GameDate: time.Now().Format("2006-01-02")}),
		Context: GameContext{
			Inning: 1,
			Half:   Top,
		},
		AwayOrder: makeBattingOrder(DefaultBattingSlots),
		HomeOrder: makeBattingOrder(DefaultBattingSlots),
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

func (b *Book) AdvanceBattingPosition() {
	_, spot := b.battingMemoryRef(b.Context.Half)
	*spot = *spot + 1
}

func (b *Book) RetreatBattingPosition() {
	_, spot := b.battingMemoryRef(b.Context.Half)
	*spot = *spot - 1
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
	b.Meta = normalizedMeta(b.Meta)
	b.AwayOrder = makeBattingOrder(b.Meta.AwaySlots)
	b.HomeOrder = makeBattingOrder(b.Meta.HomeSlots)
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
	spot = normalizeSpot(spot, len(order))
	if spot < 0 || spot >= len(order) {
		return ""
	}
	return order[spot]
}

func (b Book) BattingSequence() int {
	return b.battingSequence(b.Context.Half)
}

func (b Book) battingSequence(half Half) int {
	_, spot := b.battingMemory(half)
	if spot < 0 {
		return normalizeSpot(spot, b.battingSlots(half)) + 1
	}
	return spot + 1
}

func (b Book) BattingPosition() int {
	_, spot := b.battingMemory(b.Context.Half)
	return normalizeSpot(spot, b.battingSlots(b.Context.Half)) + 1
}

func (b Book) BattingSequenceForEntry(id string) int {
	replay := NewBook()
	for _, entry := range b.Entries {
		if entry.ID == id {
			if entry.BattingPos > 0 {
				return entry.BattingPos
			}
			return replay.battingSequence(entry.Half)
		}
		replay.RecordPlateAppearance(entry)
	}
	return b.BattingSequence()
}

func (b Book) BattingPositionForEntry(id string) int {
	for _, entry := range b.Entries {
		if entry.ID == id {
			return normalizeSpot(b.BattingSequenceForEntry(id)-1, b.battingSlots(entry.Half)) + 1
		}
	}
	return normalizeSpot(b.BattingSequenceForEntry(id)-1, b.battingSlots(b.Context.Half)) + 1
}

func CountPitches(pitches string) int {
	count := 0
	for _, pitch := range strings.TrimSpace(pitches) {
		switch pitch {
		case '+', '*', '.', '1', '2', '3', '>', 'A', 'N', 'V':
			continue
		default:
			count++
		}
	}
	return count
}

func (b Book) PitchCountForPitcher(pitcher, editingID, currentPitches string) int {
	pitcher = strings.TrimSpace(pitcher)
	if pitcher == "" {
		return 0
	}

	count := 0
	for _, entry := range b.Entries {
		if entry.ID == editingID {
			continue
		}
		if strings.TrimSpace(entry.Pitcher) != pitcher {
			continue
		}
		count += CountPitches(entry.Pitches)
	}

	return count + CountPitches(currentPitches)
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
	*order = normalizeBattingOrder(*order, b.battingSlots(entry.Half))
	slots := len(*order)
	if entry.BattingPos > 0 {
		position := entry.BattingPos - 1
		(*order)[normalizeSpot(position, slots)] = batter
		*spot = position + 1
		return
	}
	currentSpot := *spot
	if currentSpot < 0 {
		currentSpot = normalizeSpot(currentSpot, slots)
	}
	index := indexOfBatter(*order, batter)
	if index >= 0 {
		position := currentSpot - normalizeSpot(currentSpot, slots) + index
		if position < currentSpot {
			position += slots
		}
		*spot = position + 1
		return
	}

	(*order)[normalizeSpot(currentSpot, slots)] = batter
	*spot = currentSpot + 1
}

func (b Book) battingSlots(half Half) int {
	if half == Bottom {
		return normalizeBattingSlots(b.Meta.HomeSlots)
	}
	return normalizeBattingSlots(b.Meta.AwaySlots)
}

func (b Book) BattingSlots() int {
	return b.battingSlots(b.Context.Half)
}

func (b Book) BattingSlotsForHalf(half Half) int {
	return b.battingSlots(half)
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

func makeBattingOrder(slots int) []string {
	return make([]string, normalizeBattingSlots(slots))
}

func normalizeBattingOrder(order []string, slots int) []string {
	slots = normalizeBattingSlots(slots)
	if len(order) == slots {
		return order
	}
	normalized := makeBattingOrder(slots)
	copy(normalized, order)
	return normalized
}

func normalizeBattingSlots(slots int) int {
	if slots < 1 {
		return DefaultBattingSlots
	}
	return slots
}

func normalizedMeta(meta GameMeta) GameMeta {
	meta.AwaySlots = normalizeBattingSlots(meta.AwaySlots)
	meta.HomeSlots = normalizeBattingSlots(meta.HomeSlots)
	return meta
}

func normalizeSpot(spot, slots int) int {
	slots = normalizeBattingSlots(slots)
	if slots == 0 {
		return 0
	}
	spot %= slots
	if spot < 0 {
		spot += slots
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

func (e EventEntry) EventText() string {
	batterEvent := strings.TrimSpace(e.BatterEvent)
	runnerEvent := strings.TrimSpace(e.RunnerEvent)

	if batterEvent == "" {
		return runnerEvent
	}
	if runnerEvent == "" {
		return batterEvent
	}
	return batterEvent + "+" + runnerEvent
}

func (e EventEntry) Summary() string {
	summary := fmt.Sprintf("%s | batter %s | %s", strings.Title(string(e.Half)), e.Batter, e.LogEventText())
	return summary
}

func (e EventEntry) LogEventText() string {
	parts := []string{e.EventText()}
	if strings.TrimSpace(e.Advances) != "" {
		parts = append(parts, strings.TrimSpace(e.Advances))
	}
	return strings.Join(parts, ".")
}

func nextID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
