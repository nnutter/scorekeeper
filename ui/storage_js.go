//go:build js && wasm

package ui

import (
	"errors"
	"syscall/js"
)

const storageKey = "scorekeeper-book"

func loadSavedBook() (string, error) {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return "", nil
	}
	value := storage.Call("getItem", storageKey)
	if value.IsNull() || value.IsUndefined() {
		return "", nil
	}
	return value.String(), nil
}

func saveBook(raw string) error {
	storage := js.Global().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return nil
	}
	storage.Call("setItem", storageKey, raw)
	return nil
}

func copyText(raw string) error {
	navigator := js.Global().Get("navigator")
	if navigator.IsUndefined() || navigator.IsNull() {
		return errors.New("clipboard unavailable")
	}
	clipboard := navigator.Get("clipboard")
	if clipboard.IsUndefined() || clipboard.IsNull() {
		return errors.New("clipboard unavailable")
	}
	clipboard.Call("writeText", raw)
	return nil
}

func clearEntryFields(keepBatter, keepPitches bool) {
	document := js.Global().Get("document")
	if document.IsUndefined() || document.IsNull() {
		return
	}
	selectors := []string{
		`input[id^="batter-event-"]`,
		`input[id^="advances-"]`,
		`input[id^="runner-event-"]`,
		`textarea[id^="note-"]`,
	}
	if !keepBatter {
		selectors = append(selectors, `input[id^="batter-"]`)
	}
	if !keepPitches {
		selectors = append(selectors, `input[id^="pitches-"]`)
	}
	for _, selector := range selectors {
		nodes := document.Call("querySelectorAll", selector)
		for i := 0; i < nodes.Length(); i++ {
			nodes.Index(i).Set("value", "")
		}
	}
}
