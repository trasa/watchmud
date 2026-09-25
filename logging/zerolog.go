package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initialize logging to stdout, and to file as well when one is named.
// Returns the close for that file, for the caller to defer.
//
// No file is stdout only, which is what a container wants: docker (or
// kubernetes) keeps it and rotates it, and a file inside the container
// would grow until the disk filled.
func Initialize(file, level string) (func(), error) {
	if file == "" {
		return initialize(os.Stdout, level, func() {})
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return initialize(io.MultiWriter(os.Stdout, f), level, func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing log file %s: %v\n", file, err)
		}
	})
}

func initialize(wrt io.Writer, level string, closeLog func()) (func(), error) {
	if logLevel, err := zerolog.ParseLevel(level); err != nil {
		closeLog()
		return nil, err
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: wrt})
	log.Info().Msg("Logging initialized.")
	return closeLog, nil
}
