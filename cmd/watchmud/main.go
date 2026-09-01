package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/logging"
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/serverconfig"
	"github.com/trasa/watchmud/world"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "watchmud: %v\n", err)
		os.Exit(1)
	}
}

// Run holds everything main used to do, so that deferred cleanup actually
// happens. Nothing below here calls os.Exit.
func run() error {
	configPath := flag.String("config", "app.local.yaml", "location of the server configuration file")
	contentPath := flag.String("content", "", "override location of the content files")
	flag.Parse()

	// load the serverconfig.Config from YAML
	cfg, err := serverconfig.Load(*configPath)
	if err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// an explicitly passed flag overrides the config file. Defaults live in
	// the config struct, not in the flag definition.
	if *contentPath != "" {
		cfg.ContentPath = *contentPath
	}

	closeLog, err := logging.Initialize(cfg.Log.File, cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("initializing logging: %w", err)
	}
	defer closeLog()
	log.Info().Msg("Logging initialized.")
	curdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd: %w", err)
	}
	log.Info().Msgf("Current Directory: %s", curdir)
	log.Info().Msgf("Configuration Path: %s", *configPath)
	log.Info().Msgf("Content Path: %s", cfg.ContentPath)

	// canceled on SIGINT or SIGTERM. A second signal kills the process
	//outright, which is what you want if shutdown hangs.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGKILL)
	defer stop()

	// load the content: rules and world files
	contentDir := os.DirFS(cfg.ContentPath)
	content, err := loader.LoadContent(contentDir)
	if err != nil {
		return fmt.Errorf("loading content: %w", err)
	}

	w, err := world.New(content)
	if err != nil {
		return fmt.Errorf("loading world: %w", err)
	}
	gameServer := server.New(w)
	if err := gameServer.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("game server: %w", err)
	}
	return nil
}
