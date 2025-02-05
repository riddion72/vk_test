package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"main/config"
	migrate "main/pkg/postgres/migrations"

	_ "github.com/lib/pq"
)

func initPostgresClient(cfg config.Config) (*sql.DB, error) {
	options := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode)

	database, err := sql.Open("postgres", options)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		defer database.Close()
		log.Println(err)
		return nil, err
	}
	log.Println("database connection successful")
	return database, nil
}

func ConnectionDB() (*sql.DB, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := initPostgresClient(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func MigrateDB() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := initPostgresClient(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = migrate.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return nil
}
