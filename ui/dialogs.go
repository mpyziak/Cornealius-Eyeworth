package ui

import "fyne.io/fyne/v2"

const (
	dialogSchedule = "schedule"
	dialogLanguage = "language"
	dialogAbout    = "about"
	dialogHelp     = "help"
)

// One window per key. Every extra dialog is another HWND with its own GL
// context. Unlocked - menu actions and SetOnClosed are both main-goroutine.
var openDialogs = map[string]fyne.Window{}

// Callers return early when focus is true.
func focusExisting(key string) bool {
	win, ok := openDialogs[key]
	if !ok {
		return false
	}
	win.RequestFocus()
	return true
}

// Call before win.Show().
func registerDialog(key string, win fyne.Window) {
	openDialogs[key] = win
	win.SetOnClosed(func() { delete(openDialogs, key) })
}
