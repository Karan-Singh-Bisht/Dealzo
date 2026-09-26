package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/karan-singh-bisht/Dealzo-api/internal/config"
)

func main() {

	cfg := config.MustLoad()

	fmt.Println(os.Args)
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.Database_url,
	)

	if err != nil {
		log.Fatalf("migration.new : %v", err)
	}

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		//m.Down() is dangerous as it will migrate all the way down (applying all down migrations).
		//Removing all data from production db
		//To solve this we use steps
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command %s", os.Args[1])
	}
}
