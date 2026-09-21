//go:generate go run github.com/tc-hib/go-winres@latest make

package main

import (
	syslog "log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/mpyziak/cornealius-eyeworth/assets"
	"github.com/mpyziak/cornealius-eyeworth/config"
	log "github.com/mpyziak/cornealius-eyeworth/diagnostics"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/ui"
)

func main() {
	// Config sits next to the exe, not the working directory. Nothing below may
	// be fatal: under -H windowsgui there is no console to print why.
	exePath, err := os.Executable()
	if err != nil {
		syslog.Printf("warning: cannot resolve executable path (%v); using working directory", err)
		exePath = "."
	}
	exeDir := filepath.Dir(exePath)
	cfgPath := filepath.Join(exeDir, "config.json")

	if err := log.Init(exeDir); err != nil {
		syslog.Printf("warning: logging unavailable: %v", err)
	}
	defer log.Close()
	log.Info("app starting - exe=%s", exePath)

	repo := config.NewRepository(cfgPath)
	cfg, err := repo.Load()
	if err != nil {
		// Leave the bad file alone so the user can repair it.
		log.Err("cannot load config.json (%v) - starting with defaults", err)
		syslog.Printf("warning: cannot load config.json (%v); starting with defaults", err)
		cfg = config.DefaultConfig()
	}
	log.Info("config loaded - eye=%s standUp=%s lang=%s", cfg.CronExpression, cfg.StandUpCronExpression, cfg.Language)

	i18n.SetLanguage(cfg.Language)

	a := app.NewWithID("CornealiusEyeworth")
	a.Settings().SetTheme(theme.DefaultTheme())
	// Set it here. Fyne's systray onReady reads it back on its
	// own goroutine, so a later SetIcon races that read.
	a.SetIcon(assets.Icon)

	log.Info("UI starting")
	ui.Run(a, cfg, repo)
	log.Info("app exiting")
}
