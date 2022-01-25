package config

import (
	"os"
	"strconv"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logW = zerolog.ConsoleWriter{
	Out: os.Stderr,
	FormatTimestamp: func(i interface{}) string {
		parse, _ := time.Parse(time.RFC3339, i.(string))

		return parse.Format("2006-01-02 15:04:05")
	},
}

func InitializeLogger() {
	isPretty := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	if v, ok := os.LookupEnv("LOG_PRETTY"); ok {
		isPretty, _ = strconv.ParseBool(v)
	}

	if isPretty {
		log.Logger = zerolog.New(logW).With().Timestamp().Logger()
	}
}

func SetLogLevel(level string) {
	zLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Err(err).Str("component", "log").Msgf("zerolog unknown level %s", level)

		return
	}

	zerolog.SetGlobalLevel(zLevel)
}
