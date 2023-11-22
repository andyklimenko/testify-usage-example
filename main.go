package main

import (
	"github.com/andyklimenko/testify-usage-example/api"
	"github.com/andyklimenko/testify-usage-example/api/external/changelog"
	"github.com/andyklimenko/testify-usage-example/api/storage"
	"github.com/andyklimenko/testify-usage-example/api/storage/database"
	"github.com/andyklimenko/testify-usage-example/api/storage/migrations"
	"github.com/andyklimenko/testify-usage-example/config"
)

func main() {
	var cfg config.Config
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	db, err := database.DbConnect(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		panic(err)
	}

	if err := migrations.Up(db, cfg.DB.Driver); err != nil {
		panic(err)
	}

	changelogNotifySvc := changelog.New(cfg.Notify.Addr)
	srv := api.New(cfg, storage.New(db), changelogNotifySvc)
	if err := srv.Start(); err != nil {
		panic(err)
	}
}
