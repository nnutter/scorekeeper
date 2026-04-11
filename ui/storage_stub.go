//go:build !js || !wasm

package ui

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

func focusEntryField(_ string) {}

func initPullToRefresh() {}

func shouldStickToPageBottom() bool { return false }

func scrollPageToBottom() {}
