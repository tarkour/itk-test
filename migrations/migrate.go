package migrations

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"
	"github.com/tarkour/itk-test/pkg/config"
)

func Migrate() {

	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal("Config error: ", err)
	}

	db, err := sql.Open("pgx", cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		log.Fatal(err)
	}
}
