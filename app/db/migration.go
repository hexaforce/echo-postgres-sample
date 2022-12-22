package db

import (
	"embed"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

//go:embed migrations/*.sql
var sql embed.FS

func MigrateDB() (*pg.DB, error) {

	var (
		opts *pg.Options
		err  error
	)

	if os.Getenv("ENV") == "PROD" {
		opts, err = pg.ParseURL(os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, err
		}
	} else {
		opts = &pg.Options{
			Addr:     fmt.Sprintf("%s:%s", os.Getenv("POSTGRES_SERVER"), os.Getenv("POSTGRES_PORT")),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
		}
	}

	pgdb := pg.Connect(opts)
	collection := migrations.NewCollection()
	err = collection.DiscoverSQLMigrationsFromFilesystem(http.FS(sql), "migrations")
	if err != nil {
		return nil, err
	}

	_, _, err = collection.Run(pgdb, "init")
	if err != nil {
		return nil, err
	}

	oldVersion, newVersion, err := collection.Run(pgdb, "up")
	if err != nil {
		return nil, err
	}

	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	return pgdb, err

}
