package main

import (
	"log"
	"os"

	"github.com/Siddharth9101/olx-api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}

	cfg := config.MustLoad()

	 m, err := migrate.New(
        "file://migrations",
        cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	switch os.Args[1]{
	case "up":
		log.Println("running up migrations")
		if err := m.Up(); err != nil {
			log.Fatalf("migration.up: %v", err)
		}
	case "down":
		log.Println("running down migrations")
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migration.down: %v", err)
		}
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
