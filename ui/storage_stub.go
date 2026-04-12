//go:build !js || !wasm

package ui

import "github.com/nnutter/scorekeeper/internal/scorebook"

func loadSavedBook() (string, error) {
	return "", nil
}

func saveBook(_ string) error {
	return nil
}

func copyText(_ string) error {
	return nil
}

func clearEntryFields(_, _ bool) {}

func clearBatterField() {}

func syncDraftFields(_ scorebook.EventDraft) {}

func syncContextFields(_ scorebook.GameContext) {}

func focusEntryField(_ string) {}

func initPullToRefresh() {}

func shouldStickToPageBottom() bool { return false }

func scrollPageToBottom() {}
