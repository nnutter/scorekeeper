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
