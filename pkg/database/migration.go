package database

import (
	"database/sql"
	"fmt"
	"log"

	// import file base migration
	_ "github.com/golang-migrate/migrate/source/file"
	// import pq driver
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
)

func GetConnectionString(
	config *DatabaseConfig,
) string {
	var connStr = fmt.Sprintf("user=%s password=%s host=%s dbname=%s",
		config.Username, config.Password, config.Host, config.DatabaseName)
	if config.SSLMode != "" {
		connStr = fmt.Sprintf("%s sslmode=%s", connStr, config.SSLMode)
	}
	return connStr
}

func NewDatabase(
	config *DatabaseConfig,
) *sqlx.DB {
	var connStr = GetConnectionString(config)
	log.Printf("connecting to database %s@%s/%s",
		config.Username, config.Host, config.DatabaseName)
	db := sqlx.MustOpen("postgres", connStr)
	err := db.Ping()
	if err != nil {
		log.Fatal("fail to ping database", err)
	}
	log.Print("connected to database")
	return db
}

func MigrateDatabase(
	db *sql.DB, migrationPath string,
) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	migrationDir := fmt.Sprintf("file://%s", migrationPath)
	m, err := migrate.NewWithDatabaseInstance(migrationDir, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("running database migrations dir=%s\n", migrationPath)
	version, dirty, _ := m.Version()
	log.Printf("current database schema version=%d dirty=%v\n", version, dirty)

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			// not strictly necessary, but close the migration client in order to clean up
			// potentially dangling resources like locks
			closeMigrate(m)
			log.Fatal(err)
		}
		log.Print("no database migrations ran")
	} else {
		log.Print("successfully ran database migrations")
	}
}

func closeMigrate(migrate *migrate.Migrate) {
	sourceErr, databaseErr := migrate.Close()
	if sourceErr != nil {
		log.Fatal("error closing migration source", sourceErr)
	}
	if databaseErr != nil {
		log.Fatal("error closing database source", databaseErr)
	}
}
