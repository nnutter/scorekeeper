package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/nnutter/scorekeeper/internal/scorebook"
)

type Root struct {
	app.Compo

	book        scorebook.Book
	draft       scorebook.EventDraft
	message     string
	focused     string
	hasLoaded   bool
	messageKind string
	formVersion int
}

func New() *Root {
	r := &Root{}
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
	r.restore()
	ctx.Update()
}

func (r *Root) Render() app.UI {
	exportText := scorebook.ExportText(r.book)

	return app.Div().Class("page").Body(
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


func (r *Root) renderGameInfo(exportText string) app.UI {
	return app.Section().Class("panel").Body(
		app.Div().Class("game-info-row").Body(
			app.Div().Class("game-info-grid").Body(
				r.textField("Away Team", &r.book.Meta.AwayTeam, "away-team", "e.g. Yankees"),
				app.Div().Class("field").Body(
					app.Label().Text("Game Date"),
					app.Input().ID(r.fieldID("game-date")).Class("input").Type("date").Value(r.book.Meta.GameDate).
						OnInput(r.bindString(&r.book.Meta.GameDate, "game-date")).
						OnFocus(r.setFocus("game-date")),
				),
				r.textField("Home Team", &r.book.Meta.HomeTeam, "home-team", "e.g. Red Sox"),
			),
			app.Div().Class("game-info-actions").Body(
				app.Button().Class("btn").Text("New Game").OnClick(r.newGame),
				app.Button().Class("btn primary").Text("Copy").OnClick(r.copyExport),
				app.A().Class("btn").Href(scorebook.MailtoLink(r.book)).Text("Email"),
			),
		),
		app.Details().Class("export-details").Body(
			app.Summary().Class("export-summary").Text("Show Preview"),
			app.Pre().Class("panel export-box").Text(exportText),
		),
	)
}

func (r *Root) renderContext() app.UI {
	return app.Section().Class("panel context-panel").Body(
		app.Div().Class("stack").Body(
			app.Div().Class("field context-actions").Body(
				app.Label().Text(" "),
				app.Div().Class("context-action-row").Body(
					app.Button().Class("btn warm context-step").Text("-").OnClick(r.retreatHalf),
					app.Button().Class("btn warm context-step").Text("+").OnClick(r.advanceHalf),
				),
			),
			app.Div().Class("field").Body(
				app.Label().Class("field-label").Text("Inning"),
				app.Div().Class("context-chip compact").Body(
					app.Span().Text(fmt.Sprintf("%d %s", r.book.Context.Inning, strings.Title(string(r.book.Context.Half)))),
				),
			),
			r.textField("Pitcher", &r.book.Context.Pitcher, "pitcher", "e.g. 45S"),
		),
	)
}

func (r *Root) renderEntry() app.UI {
	return app.Section().Class("panel").Body(
		app.If(r.message != "", func() app.UI {
			class := "notice compact"
			if r.messageKind == "status" {
				class += " status"
			}
			return app.Div().Class(class).Text(r.message)
		}),
		app.Div().ID(r.fieldID("entry-grid")).Class(r.entryGridClass()).Body(r.renderEntryFields()...),
		app.Div().Class("actions-row").Body(
			app.Button().Class("btn primary").Text(r.saveLabel()).OnClick(r.saveEntry),
			app.If(r.draft.EditingID != "", func() app.UI {
				return app.Button().Class("btn").Text("Cancel Edit").OnClick(r.cancelEdit)
			}),
		),
	)
}

func (r *Root) statusMessage(text string) {
	r.message = text
	r.messageKind = "status"
}

func (r *Root) errorMessage(text string) {
	r.message = text
	r.messageKind = "error"
}

func (r *Root) renderEntryFields() []app.UI {
	fields := []app.UI{
		r.textField("Batter", &r.draft.Batter, "batter", "e.g. 12J"),
		r.textField("Pitches", &r.draft.Pitches, "pitches", "e.g. CBX"),
		r.textField("Batter Event", &r.draft.BatterEvent, "batter-event", "e.g. S7"),
		r.textField("Event Advances", &r.draft.Advances, "advances", "e.g. 1-3"),
		r.textField("Base-Running Event", &r.draft.RunnerEvent, "runner-event", "e.g. SB2"),
	}

	fields = append(fields, r.textAreaField("Note", &r.draft.Note, "note", "Optional note"))
	return fields
}

func (r *Root) renderKeyboard() app.UI {
	return app.Section().Class("panel keyboard-panel").Body(
		app.P().Class("meta-line").Text(r.keyboardHelpText()),
		app.Div().Class("keyboard-grid").Body(
			r.renderTokenGroup(scorebook.PitchTokenRows, "pitches"),
			r.renderTokenGroup(scorebook.BatterTokenRows, "batter-event"),
			r.renderTokenGroup(scorebook.RunnerTokenRows, "runner-event"),
			r.renderTokenGroup(scorebook.AdvanceTokenRows, "advances"),
		),
	)
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
			buttons = append(buttons, app.Button().Class("btn token").Text(t).OnClick(r.insertToken(target, t)))
		}
		uiRows = append(uiRows, app.Div().Class("keyboard-row").Body(buttons...))
	}
	return app.Div().Class("keyboard-group").Body(uiRows...)
}

func (r *Root) renderLog() app.UI {
	rows := make([]app.UI, 0, len(r.book.Entries))
	for _, entry := range r.book.Entries {
		rows = append(rows, r.renderLogEntry(entry))
	}
	if len(rows) == 0 {
		rows = append(rows, app.Div().Class("log-empty").Text("No events yet."))
	}
	return app.Section().Class("panel").Body(
		app.Div().Class("log-table").Body(
			app.Div().Class("log-row log-header").Body(
				app.Span().Text("Inning"),
				app.Span().Text("P"),
				app.Span().Text("B"),
				app.Span().Text("Pitches"),
				app.Span().Text("Event"),
				app.Span().Text("Note"),
				app.Span().Text(""),
			),
			app.Div().Class("entry-list").Body(rows...),
		),
	)
}

func (r *Root) renderLogEntry(entry scorebook.EventEntry) app.UI {
	return app.Div().Class("log-row").Body(
		app.Span().Text(shortContext(entry)),
		app.Span().Text(entry.Pitcher),
		app.Span().Text(entry.Batter),
		app.Span().Text(orDash(entry.Pitches)),
		app.Span().Text(r.logEventText(entry)),
		app.Span().Class("log-note").Text(orDash(entry.Note)),
		app.Div().Class("log-actions").Body(
			app.Button().Class("btn").Text("Edit").OnClick(r.editEntry(entry.ID)),
			app.Button().Class("btn danger").Text("Delete").OnClick(r.deleteEntry(entry.ID)),
		),
	)
}

func (r *Root) textField(label string, target *string, focusKey, placeholder string) app.UI {
	return app.Div().Class("field").Body(
		app.Label().Class("field-label").Text(label),
		app.Input().ID(r.fieldID(focusKey)).Class("input").Type("text").Value(*target).Placeholder(placeholder).
			OnInput(r.bindString(target, focusKey)).
			OnFocus(r.setFocus(focusKey)),
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
	r.clearMessage()
	r.formVersion++
	r.persist()
	ctx.Reload()
}

func (r *Root) retreatHalf(ctx app.Context, _ app.Event) {
	r.book.RetreatHalf()
	r.clearMessage()
	r.formVersion++
	r.persist()
	ctx.Reload()
}

func (r *Root) saveEntry(ctx app.Context, _ app.Event) {
	issues := scorebook.Validate(r.book.Meta, r.book.Context, r.draft)
	if len(issues) > 0 {
		r.errorMessage(issues[0])
		ctx.Update()
		return
	}

	entry := r.draft.ToEntry(r.book.Context)
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
		r.statusMessage("Event updated.")
	} else {
		r.book.Entries = append(r.book.Entries, entry)
		r.statusMessage("Event saved.")
	}
	if wasEditing {
		r.draft.Reset()
	} else if entry.Mode == scorebook.ModeRun {
		r.draft.PrepareForNextRunnerEvent()
	} else {
		r.draft.PrepareForNextPlateAppearance()
	}
	r.focused = ""
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
}

func (r *Root) cancelEdit(ctx app.Context, _ app.Event) {
	r.draft.Reset()
	r.statusMessage("Edit canceled.")
	r.formVersion++
	r.persist()
	ctx.Update()
}

func (r *Root) editEntry(id string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		for _, entry := range r.book.Entries {
			if entry.ID == id {
				r.draft.LoadFromEntry(entry)
				r.statusMessage("Editing event.")
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
		if r.draft.EditingID == id {
			r.draft.Reset()
			r.formVersion++
		}
		r.statusMessage("Event deleted.")
		r.persist()
		ctx.Update()
	}
}

func (r *Root) insertToken(target, token string) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		switch target {
		case "pitches":
			r.draft.Pitches += token
		case "batter-event":
			r.draft.BatterEvent += token
		case "advances":
			r.draft.Advances += token
		case "runner-event":
			r.draft.RunnerEvent += token
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
		r.draft.BatterEvent += token
	case "advances":
		r.draft.Advances += token
	case "runner-event":
		r.draft.RunnerEvent += token
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

func (r *Root) entryGridClass() string { return "entry-grid combined-grid" }

func (r *Root) copyExport(ctx app.Context, _ app.Event) {
	if err := copyText(scorebook.ExportText(r.book)); err != nil {
		r.errorMessage("Clipboard copy is unavailable here.")
	} else {
		r.statusMessage("Export copied.")
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
	if entry.Mode == scorebook.ModeRun {
		return entry.RunnerEvent
	}
	parts := []string{entry.BatterEvent}
	if entry.Advances != "" {
		parts = append(parts, entry.Advances)
	}
	return strings.Join(parts, " | ")
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
	book.HydratePitcherMemory()
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
	r.focused = ""
	r.statusMessage("New game started.")
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
	half := "T"
	if entry.Half == scorebook.Bottom {
		half = "B"
	}
	return fmt.Sprintf("%d%s", entry.Inning, half)
}
