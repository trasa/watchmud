package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initialize logging, both console and file specified
// returns the *os.File so that the caller can defer the close of the log
// file appropriately.
func Initialize(file, level string) (func(), error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	wrt := io.MultiWriter(os.Stdout, f)
	if logLevel, err := zerolog.ParseLevel(level); err != nil {
		return nil, err
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: wrt})
	log.Info().Msg("Logging initialized.")
	return func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing log file %s: %v\n", file, err)
		}
	}, nil
}
