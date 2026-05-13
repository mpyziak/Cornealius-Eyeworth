//go:generate go-winres make

package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/app"

	"github.com/mpyziak/cornealius-eyeworth/assets"
	"github.com/mpyziak/cornealius-eyeworth/config"
	"github.com/mpyziak/cornealius-eyeworth/i18n"
	"github.com/mpyziak/cornealius-eyeworth/ui"
)

func main() {
	// Resolve the config file path relative to the executable so it works
	// regardless of the working directory (including after fyne package).
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("cannot resolve executable path: %v", err)
	}
	cfgPath := filepath.Join(filepath.Dir(exePath), "config.json")

	repo := config.NewRepository(cfgPath)
	cfg, err := repo.Load()
	if err != nil {
		log.Fatalf("cannot load config.json: %v", err)
	}

	i18n.SetLanguage(cfg.Language)

	a := app.NewWithID("CornealiusEyeworth")
	a.SetIcon(assets.Logo)

	ui.Run(a, cfg, repo)
}
