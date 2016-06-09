package app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/postgres-ci/notifier/src/common"
	"github.com/postgres-ci/notifier/src/plugins/email"
	"github.com/postgres-ci/notifier/src/plugins/telegram"

	"os"
	"runtime"
	"time"
)

const (
	MaxOpenConns    = 5
	MaxIdleConns    = 2
	ConnMaxLifetime = time.Hour
)

func New(config common.Config) *app {

	connect, err := sqlx.Connect("postgres", config.Connect.DSN())

	if err != nil {

		log.Fatalf("Could not connect to database server: %v", err)
	}

	log.Debugf("Connect to postgresql server. DSN(%s)", config.Connect.DSN())

	connect.SetMaxOpenConns(MaxOpenConns)
	connect.SetMaxIdleConns(MaxIdleConns)
	connect.SetConnMaxLifetime(ConnMaxLifetime)

	app := app{
		config:  config,
		connect: connect,
		plugins: []plugin{
			email.New(config),
		},
	}

	if bot, err := telegram.New(config.Telegram.Token, connect); err == nil {

		app.plugins = append(app.plugins, bot)
	}

	return &app
}

type plugin interface {
	Send(common.Build) error
}

type app struct {
	config  common.Config
	connect *sqlx.DB
	plugins []plugin
	debug   bool
}

func (a *app) SetDebugMode() {

	a.debug = true
}

func (a *app) Run() {

	log.Info("Postgres-CI notifier started")
	log.Debugf("Runtime version: %s. Pid: %d", runtime.Version(), os.Getpid())

	if a.debug {

		go a.debugInfo()
	}

	a.listen()
}
