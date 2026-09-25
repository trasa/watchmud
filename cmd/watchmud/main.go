package main

import (
	"context"
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/dice"
	"github.com/trasa/watchmud/loader"
	"github.com/trasa/watchmud/logging"
	"github.com/trasa/watchmud/memstore"
	"github.com/trasa/watchmud/mongostore"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/server"
	"github.com/trasa/watchmud/serverconfig"
	"github.com/trasa/watchmud/telnet"
	"github.com/trasa/watchmud/world"
	"github.com/trasa/watchmud/writebehind"
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

	// load and verify the serverconfig.Config from YAML
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

	d, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("main: %w", err)
	}
	log.Info().Msgf("Current Directory: %s", d)
	log.Info().Msgf("Configuration Path: %s", *configPath)
	log.Info().Msgf("Content Path: %s", cfg.ContentPath)

	// canceled on SIGINT or SIGTERM. A second signal kills the process
	//outright, which is what you want if shutdown hangs.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// load the content: rules and world files
	contentDir := os.DirFS(cfg.ContentPath)
	content, err := loader.LoadContent(contentDir)
	if err != nil {
		return fmt.Errorf("loading content: %w", err)
	}

	// persistence
	store, closeStore, err := openStore(ctx, cfg)
	if err != nil {
		return fmt.Errorf("persistence: %w", err)
	}
	defer closeStore()
	saver := writebehind.New(store)
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if err := saver.Close(closeCtx); err != nil {
			log.Error().Err(err).Msg("flushing player saves")
		}
	}()

	// randomness
	var seed [32]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return fmt.Errorf("generating random seed: %w", err)
	}
	roller := dice.New(seed)

	// build the world and the game server
	w, err := world.New(content, saver, roller)
	if err != nil {
		return fmt.Errorf("loading world: %w", err)
	}
	gameServer := server.New(w, content.Catalog, saver)

	// The listeners run beside the game; if they fail -- a port in use, a
	// certificate that won't load -- the server stops with that error rather
	// than run on with nobody able to connect, looking healthy.
	ctx, stopForListener := context.WithCancelCause(ctx)
	defer stopForListener(nil)
	go func() {
		opts := telnet.Options{Addr: fmt.Sprintf("%s:%d", cfg.Telnet.Host, cfg.Telnet.Port)}
		if cfg.TLS.Port != 0 {
			opts.TLSAddr = fmt.Sprintf("%s:%d", cfg.Telnet.Host, cfg.TLS.Port)
			opts.TLSPort = cfg.TLS.Port
			opts.CertFile, opts.KeyFile = cfg.TLS.Cert, cfg.TLS.Key
		}
		if err := telnet.Listen(ctx, opts, gameServer, content.Catalog); err != nil {
			log.Error().Err(err).Msg("telnet listener")
			stopForListener(err)
		}
	}()

	// run the game server
	runErr := gameServer.Run(ctx)
	if cause := context.Cause(ctx); cause != nil && !errors.Is(cause, context.Canceled) {
		return fmt.Errorf("listener: %w", cause)
	}
	if runErr != nil && !errors.Is(runErr, context.Canceled) {
		return fmt.Errorf("game server: %w", runErr)
	}

	return nil
}

// openStore picks where characters are kept, and returns the cleanup to run on
// the way out.
//
// A configured mongo that can't be reached is a hard startup failure rather
// than a fallback to memory: a server that comes up anyway looks healthy right
// until it has quietly thrown away an evening of play. No uri at all is a
// different thing -- that is someone who asked for a throwaway server, and
// they get one, loudly.
func openStore(ctx context.Context, cfg *serverconfig.Config) (player.Store, func(), error) {
	if cfg.Mongo.Uri == "" {
		log.Warn().Msg("no mongo.uri configured: using the in-memory store, nothing will survive a restart")
		return memstore.New(), func() {}, nil
	}

	store, err := mongostore.New(ctx, cfg.Mongo.Uri, cfg.Mongo.Database)
	if err != nil {
		return nil, nil, err
	}
	log.Info().Str("uri", mongostore.RedactURI(cfg.Mongo.Uri)).Str("database", cfg.Mongo.Database).Msg("persistence: mongo")

	return store, func() {
		// ctx is cancelled by the signal that got us here, so the disconnect
		// needs a deadline of its own or it has none at all.
		closeCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if err := store.Close(closeCtx); err != nil {
			log.Error().Err(err).Msg("closing the player store")
		}
	}, nil
}
