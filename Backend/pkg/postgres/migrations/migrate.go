package migrations

import (
	"database/sql"
	"log"
	"os"
)

var sqlFile string = "pkg/postgres/migrations/sql/create_tabel.sql"

func Migrate(db *sql.DB) error {
	query, err := os.ReadFile(sqlFile)

	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migration successful")

	return nil
}
