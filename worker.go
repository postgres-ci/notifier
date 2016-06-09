package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/postgres-ci/notifier/src/common"

	"flag"
	"fmt"
	"os"
)

var (
	debug        bool
	pathToConfig string
)

const usage = `
Postgres-CI notifier

Usage:
    -c /path/to/config.yaml (if not set, then worker will use environment variables)
    -debug (enable debug mode)

Environment variables:

    APP_ADDRESS - Postgres-CI app-server address
    LOG_LEVEL   - one of: info, warning, error

    == PostgreSQL server credentials

    DB_HOST
    DB_PORT
    DB_USERNAME
    DB_PASSWORD
    DB_DATABASE 

    == SMTP server credentials

    SMTP_HOST
    SMTP_PORT
    SMTP_USERNAME
    SMTP_PASSWORD

    == Telegram Bot API credentials

    TELEGRAM_TOKEN
`

func main() {

	flag.BoolVar(&debug, "debug", false, "")
	flag.StringVar(&pathToConfig, "c", "", "")

	flag.Usage = func() {

		fmt.Println(usage)

		os.Exit(0)
	}

	flag.Parse()

	if log.IsTerminal() {

		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05 MST",
		})

	} else {

		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05 MST",
		})
	}

	config, err := common.ReadConfig(pathToConfig)

	if err != nil {

		log.Fatalf("Error reading configuration file: %v", err)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(config.LogLevel())
	}

	log.Debug(config)
}
