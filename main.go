package main

import (
	"github.com/andyklimenko/testify-usage-example/api"
	"github.com/andyklimenko/testify-usage-example/api/external/changelog"
	"github.com/andyklimenko/testify-usage-example/api/storage"
	"github.com/andyklimenko/testify-usage-example/api/storage/database"
	"github.com/andyklimenko/testify-usage-example/api/storage/migrations"
	"github.com/andyklimenko/testify-usage-example/config"
	"github.com/sirupsen/logrus"
)

func main() {
	var cfg config.Config
	if err := cfg.Load(); err != nil {
		logrus.Fatalf("loading config: %v", err)
	}

	db, err := database.DbConnect(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		logrus.Fatalf("db connect: %v", err)
	}

	if err := migrations.Up(db, cfg.DB.Driver); err != nil {
		logrus.Fatalf("applying db migrations: %v", err)
	}

	changelogNotifySvc := changelog.New(cfg.Notify.Addr)
	srv := api.New(cfg, storage.New(db), changelogNotifySvc)
	if err := srv.Start(); err != nil {
		logrus.Fatal(err)
	}
}
