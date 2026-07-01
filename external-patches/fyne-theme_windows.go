//go:build !android && !ios && !wasm && !test_web_driver

package app

import (
	"syscall"
	"time"

	"golang.org/x/sys/windows/registry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/theme"
)

const themeRegKey = `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`

// DefaultVariant returns the systems default fyne.ThemeVariant.
// Normally, you should not need this. It is extracted out of the root app package to give the
// settings app access to it.
func DefaultVariant() fyne.ThemeVariant {
	if isDark() {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, themeRegKey, registry.QUERY_VALUE)
	if err != nil { // older version of Windows will not have this key
		return false
	}
	defer k.Close()

	useLight, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil { // older version of Windows will not have this value
		return false
	}

	return useLight == 0
}

// WatchTheme calls the supplied function when the Windows dark/light theme changes.
func WatchTheme(onChanged func()) {
	// openKey is a helper to (re-)open the theme registry key with notification rights.
	openKey := func() (registry.Key, error) {
		return registry.OpenKey(registry.CURRENT_USER, themeRegKey,
			syscall.KEY_NOTIFY|registry.QUERY_VALUE)
	}

	// implementation based on an MIT-licensed Github Gist by Jeremy Black (c) 2022
	// https://gist.github.com/jerblack/1d05bbcebb50ad55c312e4d7cf1bc909
	var regNotifyChangeKeyValue *syscall.Proc
	if advapi32, err := syscall.LoadDLL("Advapi32.dll"); err == nil {
		if p, err := advapi32.FindProc("RegNotifyChangeKeyValue"); err == nil {
			regNotifyChangeKeyValue = p
		}
	}
	if regNotifyChangeKeyValue == nil {
		return
	}

	k, err := openKey()
	if err != nil {
		return // on older versions of windows the key may not exist
	}

	// debounce collapses rapid burst writes (Teams/Focus Assist writes the key
	// 3-5 times in <200 ms on call-end) into a single onChanged() call.
	const debounce = 300 * time.Millisecond
	var lastFired time.Time

	for {
		// blocks until the registry key has been changed
		ret, _, _ := regNotifyChangeKeyValue.Call(uintptr(k), 0, 0x00000001|0x00000004, 0, 0)
		if ret != 0 {
			// ERROR_KEY_DELETED (1018): Windows deleted and will recreate the key
			// during a full theme switch (Teams call-end, Focus Assist restore).
			// Poll until it reappears instead of spinning.
			k.Close()
			for {
				time.Sleep(50 * time.Millisecond)
				k, err = openKey()
				if err == nil {
					break
				}
			}
			continue
		}
		if now := time.Now(); now.Sub(lastFired) >= debounce {
			lastFired = now
			onChanged()
		}
	}
}
