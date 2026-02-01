package provider

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogging configures zerolog from the LOG_LEVEL environment variable.
// It writes to stderr using ConsoleWriter with RFC3339 timestamps.
// Default level is WARN if LOG_LEVEL is not set or invalid.
// This should be called in main() before cobra runs.
func InitLogging() {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).With().Timestamp().Logger()

	level := zerolog.WarnLevel
	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		parsed, err := zerolog.ParseLevel(envLevel)
		if err == nil {
			level = parsed
		}
	}
	zerolog.SetGlobalLevel(level)
}

// SetLevelFromFlags overrides the log level based on CLI flags.
// --debug takes precedence over -v/--verbose.
// If neither flag is set, the level defaults to WARN.
func SetLevelFromFlags(debug, verbose bool) {
	switch {
	case debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case verbose:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
}
