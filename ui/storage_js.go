//go:build js && wasm

package ui

import (
	"errors"
	"syscall/js"
)

const storageKey = "scorekeeper-book"
const pullIndicatorThreshold = 28

var pageBottomSticky bool
var pullToRefreshReady bool
var pullStartY float64
var pullActive bool
var pullMoved bool
var pullIndicator js.Value
var pullStartFunc js.Func
var pullMoveFunc js.Func
var pullEndFunc js.Func

const pullReloadThreshold = 72

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

func clearBatterField() {
	document := js.Global().Get("document")
	if document.IsUndefined() || document.IsNull() {
		return
	}
	nodes := document.Call("querySelectorAll", `input[id^="batter-"]`)
	for i := 0; i < nodes.Length(); i++ {
		nodes.Index(i).Set("value", "")
	}
}

func focusEntryField(id string) {
	window := js.Global().Get("window")
	document := js.Global().Get("document")
	if window.IsUndefined() || window.IsNull() || document.IsUndefined() || document.IsNull() {
		return
	}
	framesLeft := 4
	var callback js.Func
	callback = js.FuncOf(func(this js.Value, args []js.Value) any {
		node := document.Call("getElementById", id)
		if !node.IsUndefined() && !node.IsNull() {
			node.Call("focus")
			callback.Release()
			return nil
		}
		framesLeft--
		if framesLeft > 0 {
			window.Call("requestAnimationFrame", callback)
			return nil
		}
		callback.Release()
		return nil
	})
	window.Call("requestAnimationFrame", callback)
}

func isNearPageBottom(threshold float64) bool {
	window := js.Global().Get("window")
	if window.IsUndefined() || window.IsNull() {
		return false
	}
	document := js.Global().Get("document")
	if document.IsUndefined() || document.IsNull() {
		return false
	}
	docEl := document.Get("documentElement")
	if docEl.IsUndefined() || docEl.IsNull() {
		return false
	}
	scrollTop := window.Get("scrollY").Float()
	viewportHeight := window.Get("innerHeight").Float()
	scrollHeight := docEl.Get("scrollHeight").Float()
	return scrollHeight-(scrollTop+viewportHeight) <= threshold
}

func shouldStickToPageBottom() bool {
	if isNearPageBottom(140) {
		pageBottomSticky = true
		return true
	}
	if pageBottomSticky && isNearPageBottom(320) {
		return true
	}
	pageBottomSticky = false
	return false
}

func scrollPageToBottom() {
	window := js.Global().Get("window")
	if window.IsUndefined() || window.IsNull() {
		return
	}
	document := js.Global().Get("document")
	if document.IsUndefined() || document.IsNull() {
		return
	}
	docEl := document.Get("documentElement")
	if docEl.IsUndefined() || docEl.IsNull() {
		return
	}
	pageBottomSticky = true
	var framesLeft = 3
	var callback js.Func
	callback = js.FuncOf(func(this js.Value, args []js.Value) any {
		target := docEl.Get("scrollHeight")
		window.Call("scrollTo", map[string]any{
			"top":      target.Int(),
			"behavior": "auto",
		})
		framesLeft--
		if framesLeft > 0 {
			window.Call("requestAnimationFrame", callback)
			return nil
		}
		callback.Release()
		return nil
	})
	window.Call("requestAnimationFrame", callback)
}

func initPullToRefresh() {
	if pullToRefreshReady {
		return
	}
	window := js.Global().Get("window")
	document := js.Global().Get("document")
	if window.IsUndefined() || window.IsNull() || document.IsUndefined() || document.IsNull() {
		return
	}
	initPullIndicator(document)

	pullStartFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		touches := event.Get("touches")
		if touches.Length() == 0 {
			return nil
		}
		pullStartY = touches.Index(0).Get("clientY").Float()
		pullMoved = false
		pullActive = window.Get("scrollY").Float() <= 0
		hidePullIndicator()
		return nil
	})
	pullMoveFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		if !pullActive {
			return nil
		}
		event := args[0]
		touches := event.Get("touches")
		if touches.Length() == 0 {
			return nil
		}
		delta := touches.Index(0).Get("clientY").Float() - pullStartY
		if delta > 0 {
			pullMoved = true
			if delta >= pullIndicatorThreshold {
				setPullIndicator(delta, delta >= pullReloadThreshold)
			} else {
				hidePullIndicator()
			}
			event.Call("preventDefault")
		} else {
			hidePullIndicator()
		}
		return nil
	})
	pullEndFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
		shouldReload := false
		if pullActive && pullMoved {
			event := args[0]
			changedTouches := event.Get("changedTouches")
			if changedTouches.Length() > 0 {
				delta := changedTouches.Index(0).Get("clientY").Float() - pullStartY
				shouldReload = delta >= pullReloadThreshold
			}
		}
		hidePullIndicator()
		pullActive = false
		pullMoved = false
		if shouldReload {
			window.Get("location").Call("reload")
		}
		return nil
	})

	options := map[string]any{"passive": false}
	document.Call("addEventListener", "touchstart", pullStartFunc, options)
	document.Call("addEventListener", "touchmove", pullMoveFunc, options)
	document.Call("addEventListener", "touchend", pullEndFunc, options)
	pullToRefreshReady = true
}

func initPullIndicator(document js.Value) {
	if pullIndicator.Truthy() {
		return
	}
	pullIndicator = document.Call("createElement", "div")
	pullIndicator.Set("className", "pull-refresh-indicator")
	pullIndicator.Set("textContent", "Pull to reload")
	document.Get("body").Call("appendChild", pullIndicator)
}

func setPullIndicator(delta float64, ready bool) {
	if !pullIndicator.Truthy() {
		return
	}
	progress := delta / pullReloadThreshold
	if progress > 1 {
		progress = 1
	}
	pullIndicator.Set("className", "pull-refresh-indicator visible")
	if ready {
		pullIndicator.Get("classList").Call("add", "ready")
		pullIndicator.Set("textContent", "Release to reload")
	} else {
		pullIndicator.Get("classList").Call("remove", "ready")
		pullIndicator.Set("textContent", "Pull to reload")
	}
	opacity := 0.35 + (0.65 * progress)
	translate := -18 + int(18*progress)
	pullIndicator.Get("style").Set("opacity", opacity)
	pullIndicator.Get("style").Set("transform", "translate(-50%, "+js.ValueOf(translate).String()+"px)")
}

func hidePullIndicator() {
	if !pullIndicator.Truthy() {
		return
	}
	pullIndicator.Set("className", "pull-refresh-indicator")
	pullIndicator.Get("classList").Call("remove", "ready")
	pullIndicator.Set("textContent", "Pull to reload")
	pullIndicator.Get("style").Set("opacity", "")
	pullIndicator.Get("style").Set("transform", "")
}
