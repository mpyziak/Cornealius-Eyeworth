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
	// Resolve the config file path relative to the executable so it works
	// regardless of the working directory (including after fyne package).
	exePath, err := os.Executable()
	if err != nil {
		syslog.Fatalf("cannot resolve executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	cfgPath := filepath.Join(exeDir, "config.json")

	if err := log.Init(exeDir); err != nil {
		syslog.Printf("warning: logging unavailable: %v", err)
	}
	defer log.Close()
	log.Info("app starting — exe=%s", exePath)

	repo := config.NewRepository(cfgPath)
	cfg, err := repo.Load()
	if err != nil {
		log.Err("cannot load config.json: %v", err)
		syslog.Fatalf("cannot load config.json: %v", err)
	}
	log.Info("config loaded — eye=%s standUp=%s lang=%s", cfg.CronExpression, cfg.StandUpCronExpression, cfg.Language)

	i18n.SetLanguage(cfg.Language)

	a := app.NewWithID("CornealiusEyeworth")
	a.Settings().SetTheme(theme.DefaultTheme())
	a.SetIcon(assets.Logo)

	log.Info("UI starting")
	ui.Run(a, cfg, repo)
	log.Info("app exiting")
}
