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
