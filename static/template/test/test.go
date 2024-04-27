package test

import (
	"context"
	"log"

	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/email/emailtest"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/kernel/kerneltest"
	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/migrations"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var runner = dbtest.NewRunner(func() (*sqlx.DB, error) {
	sqlite.UseSQLite()

	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	err = migrations.Use().Up(context.Background(), db)
	if err != nil {
		return nil, err
	}

	log.Print("db loaded")

	return db, nil
})

var Run = runner.Run
var RunBenchmark = runner.RunBenchmark

var Kernel = kerneltest.NewTestKernelFactory(app.Kernel, &config.Config{
	Port:     443,
	BasePath: "https://example.test",

	Database: sqlite.NewConfig(":memory:"),
	Mail:     emailtest.NewTestMailerConfig(),
	Queue:    event.NewChannelQueueConfig(),
})
