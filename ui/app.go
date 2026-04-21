package ui

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/nnutter/scorekeeper/internal/scorebook"
)

type Root struct {
	app.Compo

	book        scorebook.Book
	draft       scorebook.EventDraft
	editContext scorebook.GameContext
	editBatter  int
	message     string
	focused     string
	hasLoaded   bool
	hasEditBase bool
	messageKind string
	messageID   int
	formVersion int
	mobileKeys  string
}

func New() *Root {
	r := &Root{mobileKeys: "pitches"}
	r.book = scorebook.NewBook()
	r.draft.Reset()
	return r
}

func (r *Root) OnMount(ctx app.Context) {
	ctx.Page().SetTitle("Scorekeeper")
	if r.hasLoaded {
		return
	}
	r.hasLoaded = true
	initPullToRefresh()
	r.restore()
	r.syncDraftBatter(false)
	ctx.Update()
}

func (r *Root) OnAppUpdate(ctx app.Context) {
	if ctx.AppUpdateAvailable() {
		ctx.Reload()
	}
}

func (r *Root) Render() app.UI {
	exportText := scorebook.MailText(r.book)

	return app.Div().Class(r.pageClass()).Body(
		app.Div().Class("stack main-stack").Body(
			r.renderGameInfo(exportText),
			r.renderLog(),
			app.Div().Class("event-layout").Body(
				r.renderContext(),
				r.renderEntry(),
			),
			r.renderKeyboard(),
		),
	)
}

func (r *Root) pageClass() string {
	return "page"
}

func (r *Root) renderGameInfo(exportText string) app.UI {
	return app.Section().Class("panel").Body(
		app.Div().Class("game-info-layout").Body(
			app.Div().Class("field game-away").Body(
				app.Label().Class("field-label").Text("Away Team"),
				app.Input().ID(r.fieldID("away-team")).Class("input").Type("text").Value(r.book.Meta.AwayTeam).Placeholder("Away Team").
					OnInput(r.bindString(&r.book.Meta.AwayTeam, "away-team")).
					OnFocus(r.setFocus("away-team")),
			),
			app.Div().Class("field game-date").Body(
				app.Label().Class("field-label").Text("Game Date"),
				app.Input().ID(r.fieldID("game-date")).Class("input").Type("date").Value(r.book.Meta.GameDate).
					OnInput(r.bindString(&r.book.Meta.GameDate, "game-date")).
					OnFocus(r.setFocus("game-date")),
			),
			app.Div().Class("field game-home").Body(
				app.Label().Class("field-label").Text("Home Team"),
				app.Input().ID(r.fieldID("home-team")).Class("input").Type("text").Value(r.book.Meta.HomeTeam).Placeholder("Home Team").
					OnInput(r.bindString(&r.book.Meta.HomeTeam, "home-team")).
					OnFocus(r.setFocus("home-team")),
			),
			r.iconButton("btn danger game-new", "/web/icon-clear-all.svg", "New Game").OnClick(r.newGame),
			r.iconButton("btn game-copy", "/web/icon-copy.svg", "Copy").OnClick(r.copyExport),
			r.iconLink("btn game-email", "/web/icon-email.svg", "Email", scorebook.MailtoLink(r.book)),
		),
		app.Details().Class("export-details").Body(
			app.Summary().Class("export-summary").Text("Show Preview"),
			app.Pre().Class("panel export-box").Text(exportText),
		),
	)
}

func (r *Root) renderContext() app.UI {
	return app.Section().Class("panel context-panel").Body(
		app.Div().Class("stack context-layout").Body(
			app.Div().Class("context-band context-band-top").Body(
				app.Div().Class("field context-actions").Body(
					app.Label().Text(" "),
					app.Div().Class("context-action-row").Body(
						app.Button().Class("btn warm context-step").Text("-").OnClick(r.retreatHalf),
						app.Button().Class("btn warm context-step").Text("+").OnClick(r.advanceHalf),
					),
				),
				app.Div().Class("field context-inning").Body(
					app.Label().Class("field-label").Text("Inning"),
					app.Div().Class("context-chip compact").Body(
						app.Span().Text(fmt.Sprintf("%s%d", string(r.book.Context.Half), r.book.Context.Inning)),
					),
				),
				app.Div().Class("field context-pitcher").Body(
					app.Label().Class("field-label").Text("Pitcher"),
					app.Input().ID(r.fieldID("pitcher")).Class("input").Type("text").Value(r.book.Context.Pitcher).Placeholder("45S").
						Attr("autocapitalize", "characters").
						Spellcheck(false).
						OnInput(r.bindString(&r.book.Context.Pitcher, "pitcher")).
						OnFocus(r.setFocus("pitcher")),
				),
				app.Div().Class("context-pitch-count").Text(r.pitchCountLabel()),
			),
			app.Div().Class("context-band context-band-bottom").Body(
				app.Div().Class("field context-batting").Body(
					app.Div().Class("context-batting-row").Body(
						app.Button().Class("btn navy context-step batting-step").Text("-").OnClick(r.retreatBatter),
						app.Button().Class("btn navy context-step batting-step").Text("+").OnClick(r.advanceBatter),
					),
				),
			),
		),
	)
}

func (r *Root) pitchCountLabel() string {
	count := r.book.PitchCountForPitcher(r.book.Context.Pitcher, r.draft.EditingID, r.draft.Pitches)
	return fmt.Sprintf("P: %d", count)
}

func (r *Root) renderEntry() app.UI {
	return app.Section().Class("panel").Body(
		app.Div().ID(r.fieldID("entry-grid")).Class(r.entryGridClass()).Body(r.renderEntryFields()...),
		app.Div().Class("actions-row").Body(
			r.saveIconButton().OnClick(r.saveEntry),
			app.If(r.message != "", func() app.UI {
				class := "notice compact action-notice"
				if r.messageKind == "status" {
					class += " status"
				}
				return app.Div().Class(class).Text(r.message)
			}),
			app.If(r.draft.EditingID != "", func() app.UI {
				return r.iconButton("btn", "/web/icon-cancel.svg", "Cancel Edit").OnClick(r.cancelEdit)
			}),
		),
	)
}

func (r *Root) statusMessage(ctx app.Context, text string) {
	r.message = text
	r.messageKind = "status"
	r.messageID++
	messageID := r.messageID
	ctx.After(2*time.Second, func(ctx app.Context) {
		if r.messageKind != "status" || r.messageID != messageID {
			return
		}
		r.clearMessage()
		ctx.Update()
	})
}

func (r *Root) errorMessage(text string) {
	r.message = text
	r.messageKind = "error"
	r.messageID++
}

func (r *Root) renderEntryFields() []app.UI {
	fields := []app.UI{
		r.textField(r.batterLabel(), &r.draft.Batter, "batter", "12J"),
		r.textField(r.pitchesLabel(), &r.draft.Pitches, "pitches", ""),
		r.textField("Batter Event", &r.draft.BatterEvent, "batter-event", ""),
		r.textField("Event Advances", &r.draft.Advances, "advances", ""),
		r.textField("Base-Running Event", &r.draft.RunnerEvent, "runner-event", ""),
	}

	fields = append(fields, r.textAreaField("Note", &r.draft.Note, "note", "Optional note"))
	return fields
}

func (r *Root) batterLabel() string {
	position := r.book.BattingPosition()
	if r.draft.EditingID != "" {
		position = r.editBatter
		if position < 1 || position > scorebook.BattingSlots {
			position = r.book.BattingPositionForEntry(r.draft.EditingID)
		}
	}
	return fmt.Sprintf("Batting %s", ordinal(position))
}

func ordinal(n int) string {
	if n%100 >= 11 && n%100 <= 13 {
		return fmt.Sprintf("%dth", n)
	}

	suffix := "th"
	switch n % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}

	return fmt.Sprintf("%d%s", n, suffix)
}

func (r *Root) pitchesLabel() string {
	balls, strikes := countBallsStrikes(r.draft.Pitches)
	return fmt.Sprintf("Pitches (%d-%d)", balls, strikes)
}

func countBallsStrikes(pitches string) (balls, strikes int) {
	for _, r := range pitches {
		switch r {
		case 'B', 'I', 'P', 'V':
			balls++
		case 'A', 'C', 'K', 'M', 'Q', 'S':
			strikes++
		case 'F', 'L', 'R':
			if strikes < 2 {
				strikes++
			}
		}
	}
	return
}

func (r *Root) renderKeyboard() app.UI {
	groups := r.keyboardGroups()
	desktopGroups := make([]app.UI, 0, len(groups))
	for _, group := range groups {
		desktopGroups = append(desktopGroups, r.renderTokenGroup(group.Rows, group.Target))
	}

	return app.Section().Class("panel keyboard-panel").Body(
		app.P().Class("meta-line").Text(r.keyboardHelpText()),
		app.Div().Class("keyboard-grid").Body(
			desktopGroups...,
		),
		app.Div().Class("keyboard-mobile").Body(
			app.Div().Class("keyboard-mobile-rail").Body(
				r.keyboardSwitch(groups[0]),
				r.keyboardSwitch(groups[2]),
			),
			app.Div().Class("keyboard-mobile-main").Body(
				r.renderMobileTokenGroups(groups)...,
			),
			app.Div().Class("keyboard-mobile-rail").Body(
				r.keyboardSwitch(groups[1]),
				r.keyboardSwitch(groups[3]),
			),
		),
	)
}

func (r *Root) renderMobileTokenGroups(groups []keyboardGroup) []app.UI {
	panes := make([]app.UI, 0, len(groups))
	for _, group := range groups {
		class := "keyboard-mobile-pane"
		if r.mobileKeys == group.Key {
			class += " active"
		}
		panes = append(panes, app.Div().Class(class).Body(
			r.renderTokenGroup(group.Rows, group.Target),
		))
	}
	return panes
}

type keyboardGroup struct {
	Key    string
	Label  string
	Target string
	Rows   [][]string
}

func (r *Root) keyboardGroups() []keyboardGroup {
	return []keyboardGroup{
		{Key: "pitches", Label: "P", Target: "pitches", Rows: scorebook.PitchTokenRows},
		{Key: "batter-event", Label: "B", Target: "batter-event", Rows: scorebook.BatterTokenRows},
		{Key: "runner-event", Label: "R", Target: "runner-event", Rows: scorebook.RunnerTokenRows},
		{Key: "advances", Label: "A", Target: "advances", Rows: scorebook.AdvanceTokenRows},
	}
}

func (r *Root) keyboardSwitch(group keyboardGroup) app.UI {
	class := "btn keyboard-switch"
	if r.mobileKeys == group.Key {
		class += " active"
	}
	return app.Button().Class(class).Text(group.Label).OnClick(r.setKeyboardGroup(group.Key))
}

func (r *Root) setKeyboardGroup(key string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		r.mobileKeys = key
		ctx.Update()
	}
}

func (r *Root) keyboardHelpText() string {
	return "Pitch buttons fill Pitches, play buttons fill Batter Event, advance buttons fill Event Advances, and base-running buttons fill Base-Running Event."
}

func (r *Root) renderTokenGroup(rows [][]string, target string) app.UI {
	uiRows := make([]app.UI, 0, len(rows))
	for _, row := range rows {
		buttons := make([]app.UI, 0, len(row))
		for _, token := range row {
			t := token
			buttons = append(buttons, app.Button().Class("btn token "+r.tokenToneClass(target, t)).Text(t).OnClick(r.insertToken(target, t)))
		}
		uiRows = append(uiRows, app.Div().Class("keyboard-row").Body(buttons...))
	}
	return app.Div().Class("keyboard-group").Body(uiRows...)
}

func (r *Root) tokenToneClass(target, token string) string {
	switch target {
	case "pitches":
		switch token {
		case "B", "H", "I", "V":
			return "token-good"
		case "A", "C", "F", "K", "L", "M", "O", "P", "Q", "R", "S", "T", "U", "X", "Y":
			return "token-bad"
		}
	case "batter-event":
		switch token {
		case "S", "D", "T", "HR", "DGR", "W", "IW", "HP", "E", "FC", "FLE", "SF", "SH":
			return "token-good"
		case "K", "GDP", "LDP", "FO":
			return "token-bad"
		}
	case "advances":
		switch {
		case strings.Contains(token, "X"):
			return "token-bad"
		case strings.Contains(token, "-") || token == "E" || token == "TH" || token == "RBI":
			return "token-good"
		}
	case "runner-event":
		switch token {
		case "SB2", "SB3", "SBH", "WP", "BK", "DI", "OA", "PB":
			return "token-good"
		case "CS2", "CS3", "CSH", "PO1", "PO2", "PO3", "POCS2", "POCS3", "POCSH":
			return "token-bad"
		}
	}

	return "token-neutral"
}

func (r *Root) renderLog() app.UI {
	entries := sortedLogEntries(r.book)
	rows := make([]app.UI, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, r.renderLogEntry(entry))
	}
	if len(rows) == 0 {
		rows = append(rows, app.Div().Class("log-empty").Text("No events yet."))
	}
	return app.Section().Class("panel").Body(
		app.Div().Class("log-table").Body(
			app.Div().Class("log-row log-header").Body(
				app.Span().Body(
					app.Span().Class("desktop-only").Text("Inning"),
					app.Span().Class("mobile-only").Text("Inn"),
				),
				app.Span().Text("P"),
				app.Span().Text("B"),
				app.Span().Text("Pitches"),
				app.Span().Text("Event"),
				app.Span().Text(""),
			),
			app.Div().Class("entry-list").Body(rows...),
		),
	)
}

func sortedLogEntries(book scorebook.Book) []scorebook.EventEntry {
	sorted := slices.Clone(book.Entries)
	legacyPositions := legacyBattingPositions(book.Entries)
	currentPositions := make(map[string]int, len(book.Entries))
	for _, entry := range book.Entries {
		currentPositions[entry.ID] = displayBattingPosition(entry, legacyPositions)
	}
	slices.SortStableFunc(sorted, func(a, b scorebook.EventEntry) int {
		if a.Inning != b.Inning {
			return a.Inning - b.Inning
		}
		if a.Half != b.Half {
			return halfSortRank(a.Half) - halfSortRank(b.Half)
		}
		if currentPositions[a.ID] != currentPositions[b.ID] {
			return currentPositions[a.ID] - currentPositions[b.ID]
		}
		if hasStoredBattingPosition(a) != hasStoredBattingPosition(b) && legacyPositions[a.ID] != legacyPositions[b.ID] {
			return legacyPositions[b.ID] - legacyPositions[a.ID]
		}
		return 0
	})
	return sorted
}

func hasStoredBattingPosition(entry scorebook.EventEntry) bool {
	return entry.BattingPos >= 1 && entry.BattingPos <= scorebook.BattingSlots
}

func displayBattingPosition(entry scorebook.EventEntry, legacyPositions map[string]int) int {
	if hasStoredBattingPosition(entry) {
		return entry.BattingPos
	}
	return legacyPositions[entry.ID]
}

func legacyBattingPositions(entries []scorebook.EventEntry) map[string]int {
	positions := make(map[string]int, len(entries))
	replay := scorebook.NewBook()
	for _, entry := range entries {
		positions[entry.ID] = replay.BattingPosition()
		legacyEntry := entry
		legacyEntry.BattingPos = 0
		replay.RecordPlateAppearance(legacyEntry)
	}
	return positions
}

func halfSortRank(half scorebook.Half) int {
	if half == scorebook.Top {
		return 0
	}
	return 1
}

func (r *Root) renderLogEntry(entry scorebook.EventEntry) app.UI {
	children := []app.UI{
		app.Div().Class("log-row").Body(
			app.Span().Text(shortContext(entry)),
			app.Span().Text(entry.Pitcher),
			app.Span().Text(entry.Batter),
			app.Span().Text(orDash(entry.Pitches)),
			app.Span().Text(r.logEventText(entry)),
			app.Div().Class("log-actions").Body(
				app.Button().Class("btn icon-btn").Attr("aria-label", "Edit event").Attr("title", "Edit event").Body(
					app.Img().Src("/web/icon-edit.svg").Alt(""),
				).OnClick(r.editEntry(entry.ID)),
				app.Button().Class("btn danger icon-btn").Attr("aria-label", "Delete event").Attr("title", "Delete event").Body(
					app.Img().Src("/web/icon-delete.svg").Alt(""),
				).OnClick(r.deleteEntry(entry.ID)),
			),
		),
	}
	if strings.TrimSpace(entry.Note) != "" {
		children = append(children, app.Div().Class("log-note-row").Body(
			app.Span().Class("log-note-label"),
			app.Span().Class("log-note").Text(entry.Note),
		))
	}
	return app.Div().Class("log-entry").Body(
		children...,
	)
}

func (r *Root) textField(label string, target *string, focusKey, placeholder string) app.UI {
	input := app.Input().ID(r.fieldID(focusKey)).Class("input").Type("text").Value(*target).Placeholder(placeholder).
		OnInput(r.bindString(target, focusKey)).
		OnFocus(r.setFocus(focusKey))
	if focusKey == "pitcher" || focusKey == "batter" {
		input = input.
			Attr("autocapitalize", "characters").
			Spellcheck(false)
	}

	return app.Div().Class("field").Body(
		app.Label().Class("field-label").Text(label),
		input,
	)
}

func (r *Root) textAreaField(label string, target *string, focusKey, placeholder string) app.UI {
	return app.Div().Class("field").Body(
		app.Label().Class("field-label").Text(label),
		app.Textarea().ID(r.fieldID(focusKey)).Class("textarea").Text(*target).Placeholder(placeholder).
			OnInput(r.bindString(target, focusKey)).
			OnFocus(r.setFocus(focusKey)),
	)
}

func (r *Root) iconButton(className, src, label string) app.HTMLButton {
	return app.Button().Class(className+" action-icon-btn").Attr("aria-label", label).Attr("title", label).Body(
		app.Img().Class("action-icon").Src(src).Alt(""),
	)
}

func (r *Root) iconLink(className, src, label, href string) app.HTMLA {
	return app.A().Class(className+" action-icon-btn").Href(href).Attr("aria-label", label).Attr("title", label).Body(
		app.Img().Class("action-icon").Src(src).Alt(""),
	)
}

func (r *Root) saveIconButton() app.HTMLButton {
	if r.draft.EditingID != "" {
		return r.iconButton("btn save-action-btn", "/web/icon-approve.svg", "Update Event")
	}
	return r.iconButton("btn save-action-btn", "/web/icon-add.svg", "Save Event")
}

func (r *Root) bindString(target *string, focusKey string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		*target = ctx.JSSrc().Get("value").String()
		if focusKey == "pitcher" {
			r.book.SyncPitcherMemory()
		}
		r.focused = focusKey
		if r.messageKind == "error" {
			r.clearMessage()
		}
		r.persist()
		ctx.Update()
	}
}

func (r *Root) setFocus(focusKey string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		r.focused = focusKey
		ctx.Update()
	}
}

func (r *Root) advanceHalf(ctx app.Context, _ app.Event) {
	r.book.AdvanceHalf()
	r.handleHalfChange()
	r.formVersion++
	r.persist()
	ctx.Update()
	syncContextFields(r.book.Context)
	if r.draft.EditingID == "" {
		syncDraftFields(r.draft)
	}
}

func (r *Root) retreatHalf(ctx app.Context, _ app.Event) {
	r.book.RetreatHalf()
	r.handleHalfChange()
	r.formVersion++
	r.persist()
	ctx.Update()
	syncContextFields(r.book.Context)
	if r.draft.EditingID == "" {
		syncDraftFields(r.draft)
	}
}

func (r *Root) advanceBatter(ctx app.Context, _ app.Event) {
	r.stepBatter(1)
	if strings.TrimSpace(r.draft.Batter) == "" {
		clearBatterField()
	}
	r.clearMessage()
	r.persist()
	ctx.Update()
}

func (r *Root) retreatBatter(ctx app.Context, _ app.Event) {
	r.stepBatter(-1)
	if strings.TrimSpace(r.draft.Batter) == "" {
		clearBatterField()
	}
	r.clearMessage()
	r.persist()
	ctx.Update()
}

func (r *Root) saveEntry(ctx app.Context, _ app.Event) {
	stickToBottom := shouldStickToPageBottom()

	issues := scorebook.Validate(r.book.Meta, r.book.Context, r.draft)
	if len(issues) > 0 {
		r.errorMessage(issues[0])
		ctx.Update()
		return
	}

	entry := r.draft.ToEntry(r.book.Context)
	entry.BattingPos = r.currentEntryBattingPosition()
	wasEditing := r.draft.EditingID != ""
	if r.draft.EditingID != "" {
		entry.ID = r.draft.EditingID
		for i := range r.book.Entries {
			if r.book.Entries[i].ID == r.draft.EditingID {
				entry.CreatedAt = r.book.Entries[i].CreatedAt
				r.book.Entries[i] = entry
				break
			}
		}
		r.statusMessage(ctx, "Event updated.")
	} else {
		r.book.Entries = append(r.book.Entries, entry)
		r.statusMessage(ctx, "Event saved.")
	}
	r.book.HydrateMemory()
	if wasEditing {
		r.restoreEditContext()
		r.draft.Reset()
	} else if entry.Mode == scorebook.ModeRun {
		r.draft.PrepareForNextRunnerEvent()
	} else {
		r.draft.PrepareForNextPlateAppearance()
		r.syncDraftBatter(true)
	}
	if wasEditing {
		r.syncDraftBatter(true)
	}
	shouldFocusBatter := strings.TrimSpace(r.book.RememberedBatter()) == ""
	r.focused = ""
	r.mobileKeys = "pitches"
	r.formVersion++
	r.persist()
	if wasEditing {
		clearEntryFields(false, false)
	} else if entry.Mode == scorebook.ModeRun {
		clearEntryFields(true, true)
	} else {
		clearEntryFields(false, false)
	}
	ctx.Update()
	if shouldFocusBatter {
		r.focused = "batter"
		focusEntryField(r.fieldID("batter"))
	}
	if stickToBottom {
		scrollPageToBottom()
	}
}

func (r *Root) cancelEdit(ctx app.Context, _ app.Event) {
	r.restoreEditContext()
	r.draft.Reset()
	r.syncDraftBatter(true)
	r.statusMessage(ctx, "Edit canceled.")
	r.formVersion++
	r.persist()
	clearEntryFields(false, false)
	ctx.Update()
}

func (r *Root) editEntry(id string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		for _, entry := range r.book.Entries {
			if entry.ID == id {
				if !r.hasEditBase {
					r.editContext = r.book.Context
					r.hasEditBase = true
				}
				r.book.Context = scorebook.GameContext{
					Inning:  entry.Inning,
					Half:    entry.Half,
					Pitcher: entry.Pitcher,
				}
				r.editBatter = entry.BattingPos
				if r.editBatter < 1 || r.editBatter > scorebook.BattingSlots {
					r.editBatter = r.book.BattingPositionForEntry(entry.ID)
				}
				r.draft.LoadFromEntry(entry)
				r.statusMessage(ctx, "Editing event.")
				r.formVersion++
				r.persist()
				ctx.Update()
				return
			}
		}
	}
}

func (r *Root) deleteEntry(id string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		entries := r.book.Entries[:0]
		for _, entry := range r.book.Entries {
			if entry.ID != id {
				entries = append(entries, entry)
			}
		}
		r.book.Entries = entries
		r.book.HydrateMemory()
		if r.draft.EditingID == id {
			r.restoreEditContext()
			r.draft.Reset()
			r.syncDraftBatter(true)
			r.formVersion++
		} else if r.draft.EditingID == "" {
			r.syncDraftBatter(false)
		}
		r.statusMessage(ctx, "Event deleted.")
		r.persist()
		ctx.Update()
	}
}

func (r *Root) restoreEditContext() {
	if !r.hasEditBase {
		return
	}
	r.book.Context = r.editContext
	r.editBatter = 0
	r.hasEditBase = false
}

func (r *Root) resetEntryDraftForContextChange() {
	r.draft.Reset()
	r.editBatter = 0
	r.hasEditBase = false
	r.focused = ""
	r.mobileKeys = "pitches"
	r.syncDraftBatter(true)
}

func (r *Root) stepBatter(delta int) {
	if r.draft.EditingID != "" {
		if r.editBatter < 1 || r.editBatter > scorebook.BattingSlots {
			r.editBatter = r.book.BattingPositionForEntry(r.draft.EditingID)
		}
		r.editBatter = wrapBattingPosition(r.editBatter + delta)
		r.syncEditingDraftBatter()
		return
	}
	if delta > 0 {
		r.book.AdvanceBattingPosition()
	} else {
		r.book.RetreatBattingPosition()
	}
	r.syncDraftBatter(true)
}

func (r *Root) currentEntryBattingPosition() int {
	if r.draft.EditingID != "" && r.editBatter >= 1 && r.editBatter <= scorebook.BattingSlots {
		return r.editBatter
	}
	return r.book.BattingPosition()
}

func wrapBattingPosition(position int) int {
	position = (position - 1) % scorebook.BattingSlots
	if position < 0 {
		position += scorebook.BattingSlots
	}
	return position + 1
}

func (r *Root) syncEditingDraftBatter() {
	r.draft.Batter = r.rememberedBatterAt(r.book.Context.Half, r.editBatter)
}

func (r *Root) rememberedBatterAt(half scorebook.Half, position int) string {
	if position < 1 || position > scorebook.BattingSlots {
		return ""
	}
	var order []string
	if half == scorebook.Bottom {
		order = r.book.HomeOrder
	} else {
		order = r.book.AwayOrder
	}
	index := position - 1
	if index >= len(order) {
		return ""
	}
	return strings.TrimSpace(order[index])
}

func (r *Root) handleHalfChange() {
	if r.draft.EditingID != "" {
		r.clearMessage()
		return
	}
	r.resetEntryDraftForContextChange()
	r.clearMessage()
}

func (r *Root) insertToken(target, token string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		switch target {
		case "pitches":
			r.draft.Pitches += token
		case "batter-event":
			r.appendBatterEventToken(token)
		case "advances":
			r.appendAdvanceToken(token)
		case "runner-event":
			r.appendRunnerEventToken(token)
		case "note":
			r.draft.Note += token
		case "pitcher":
			r.book.Context.Pitcher += token
			r.book.SyncPitcherMemory()
		case "batter":
			r.draft.Batter += token
		default:
			r.applyFocusedFallback(token)
		}
		r.persist()
		ctx.Update()
	}
}

func (r *Root) applyFocusedFallback(token string) {
	switch r.focused {
	case "pitches":
		r.draft.Pitches += token
	case "batter-event":
		r.appendBatterEventToken(token)
	case "advances":
		r.appendAdvanceToken(token)
	case "runner-event":
		r.appendRunnerEventToken(token)
	case "note":
		r.draft.Note += token
	case "pitcher":
		r.book.Context.Pitcher += token
		r.book.SyncPitcherMemory()
	case "batter":
		r.draft.Batter += token
	default:
		r.draft.Pitches += token
	}
}

func (r *Root) appendAdvanceToken(token string) {
	r.draft.Advances = upsertAdvanceToken(r.draft.Advances, token)
	r.draft.Advances = sortLeadRunnerFirst(r.draft.Advances, advanceSortRank)
}

func (r *Root) appendBatterEventToken(token string) {
	if isExclusiveBatterEventToken(token) {
		r.draft.BatterEvent = token + batterEventSuffix(r.draft.BatterEvent)
		return
	}
	r.draft.BatterEvent += formatBatterEventToken(r.draft.BatterEvent, token)
}

func (r *Root) appendRunnerEventToken(token string) {
	if strings.TrimSpace(r.draft.RunnerEvent) == "" {
		r.draft.RunnerEvent = token
		r.draft.RunnerEvent = sortLeadRunnerFirst(r.draft.RunnerEvent, runnerEventSortRank)
		return
	}
	if strings.HasSuffix(strings.TrimSpace(r.draft.RunnerEvent), ";") {
		r.draft.RunnerEvent += token
		r.draft.RunnerEvent = sortLeadRunnerFirst(r.draft.RunnerEvent, runnerEventSortRank)
		return
	}
	r.draft.RunnerEvent += ";" + token
	r.draft.RunnerEvent = sortLeadRunnerFirst(r.draft.RunnerEvent, runnerEventSortRank)
}

func sortLeadRunnerFirst(value string, rank func(string) int) string {
	parts := strings.Split(value, ";")
	slices.SortStableFunc(parts, func(a, b string) int {
		return rank(a) - rank(b)
	})
	return strings.Join(parts, ";")
}

func upsertAdvanceToken(value, token string) string {
	runner, ok := advanceRunnerKey(token)
	if !ok {
		if strings.TrimSpace(value) == "" {
			return token
		}
		if strings.HasSuffix(strings.TrimSpace(value), ";") {
			return value + token
		}
		return value + ";" + token
	}

	parts := strings.Split(value, ";")
	for i, part := range parts {
		if partRunner, ok := advanceRunnerKey(part); ok && partRunner == runner {
			parts[i] = token
			return strings.Join(parts, ";")
		}
	}

	if strings.TrimSpace(value) == "" {
		return token
	}
	if strings.HasSuffix(strings.TrimSpace(value), ";") {
		return value + token
	}
	return value + ";" + token
}

func advanceRunnerKey(token string) (byte, bool) {
	token = strings.TrimSpace(token)
	if token == "" {
		return 0, false
	}

	switch token[0] {
	case '3', '2', '1', 'B':
		return token[0], true
	default:
		return 0, false
	}
}

func advanceSortRank(token string) int {
	runner, ok := advanceRunnerKey(token)
	if !ok {
		return 99
	}

	switch runner {
	case '3':
		return 0
	case '2':
		return 1
	case '1':
		return 2
	case 'B':
		return 3
	default:
		return 4
	}
}

func runnerEventSortRank(token string) int {
	token = strings.TrimSpace(token)
	if token == "" {
		return 99
	}

	switch {
	case strings.HasSuffix(token, "H") || token == "PO3" || token == "POCSH":
		return 0
	case strings.HasSuffix(token, "3") || token == "PO2" || token == "POCS3":
		return 1
	case strings.HasSuffix(token, "2") || token == "PO1" || token == "POCS2":
		return 2
	default:
		return 3
	}
}

func isExclusiveBatterEventToken(token string) bool {
	switch token {
	case "K", "S", "D", "T", "HR", "DGR", "W", "IW", "HP", "E", "FC", "FLE":
		return true
	default:
		return false
	}
}

func batterEventSuffix(event string) string {
	for i, r := range event {
		if r == '(' || r == '/' {
			return event[i:]
		}
	}
	return ""
}

func formatBatterEventToken(current, token string) string {
	if !isSlashPrefixedBatterEventToken(token) {
		return token
	}
	if strings.HasSuffix(current, "/") {
		return token
	}
	return "/" + token
}

func isSlashPrefixedBatterEventToken(token string) bool {
	switch token {
	case "SF", "SH", "GDP", "LDP", "FO", "G", "L", "P", "F":
		return true
	default:
		return false
	}
}

func (r *Root) entryGridClass() string { return "entry-grid combined-grid" }

func (r *Root) copyExport(ctx app.Context, _ app.Event) {
	if err := copyText(scorebook.ExportText(r.book)); err != nil {
		r.errorMessage("Clipboard copy is unavailable here.")
	} else {
		r.statusMessage(ctx, "Export copied.")
	}
	ctx.Update()
}

func (r *Root) saveLabel() string {
	if r.draft.EditingID != "" {
		return "Update Event"
	}
	return "Save Event"
}

func (r *Root) logEventText(entry scorebook.EventEntry) string {
	return entry.LogEventText()
}

func (r *Root) restore() {
	raw, err := loadSavedBook()
	if err != nil || raw == "" {
		return
	}
	var book scorebook.Book
	if err := json.Unmarshal([]byte(raw), &book); err != nil {
		return
	}
	if book.Context.Inning < 1 {
		book.Context.Inning = 1
	}
	if book.Context.Half == "" {
		book.Context.Half = scorebook.Top
	}
	book.HydrateMemory()
	r.book = book
}

func (r *Root) persist() {
	raw, err := json.Marshal(r.book)
	if err != nil {
		return
	}
	_ = saveBook(string(raw))
}

func (r *Root) clearMessage() {
	r.message = ""
	r.messageKind = ""
}

func (r *Root) newGame(ctx app.Context, _ app.Event) {
	r.book = scorebook.NewBook()
	r.draft.Reset()
	r.syncDraftBatter(true)
	r.focused = ""
	r.statusMessage(ctx, "New game started.")
	r.formVersion++
	r.persist()
	ctx.Reload()
}

func (r *Root) fieldID(name string) string {
	return fmt.Sprintf("%s-%d", name, r.formVersion)
}

func orFallback(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func orDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
}

func shortContext(entry scorebook.EventEntry) string {
	return fmt.Sprintf("%s%d", string(entry.Half), entry.Inning)
}

func (r *Root) syncDraftBatter(force bool) {
	if r.draft.EditingID != "" {
		return
	}
	if !force && strings.TrimSpace(r.draft.Batter) != "" {
		return
	}
	remembered := strings.TrimSpace(r.book.RememberedBatter())
	if remembered == "" {
		r.draft.Batter = ""
		return
	}
	r.draft.Batter = remembered
}
