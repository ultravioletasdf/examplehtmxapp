package utils

import (
	"database/sql"
	"errors"
	sqlc "examplehtmxapp/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
	_ "modernc.org/sqlite"
)

type Database struct {
	dir       string
	connector *libsql.Connector
	db        *sql.DB
	Executor  *sqlc.Queries
}

func ConnectDatabase(cfg *Config) Database {
	if cfg.DatabaseTursoToken == "" || cfg.DatabaseTursoUrl.Host == "" {
		fmt.Printf("Connecting to local database %s\n", cfg.DatabaseFilename)
		return LocalDatabase(cfg)
	}
	fmt.Println("Connecting to turso database")
	return TursoDatabase(cfg)
}

// Cleanup function, deletes temporary directory, closes sql and connector
func (database Database) Close() error {
	errs := make([]error, 3)
	if database.dir != "" {
		errs[0] = os.RemoveAll(database.dir)
	}
	if database.connector != nil {
		errs[1] = database.connector.Close()
	}
	if database.db != nil {
		errs[2] = database.db.Close()
	}
	return errors.Join(errs...)
}

func LocalDatabase(cfg *Config) Database {
	db, err := sql.Open("sqlite", cfg.DatabaseFilename)
	if err != nil {
		log.Fatalf("Failed to connect to local database: %v", err)
	}

	return Database{db: db, Executor: sqlc.New(db)}
}

func TursoDatabase(cfg *Config) (database Database) {
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		log.Fatalf("Failed to connect to turso database: error creating temporary directory: %v", err)
	}
	database.dir = dir

	dbFile := filepath.Join(dir, cfg.DatabaseFilename)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbFile, cfg.DatabaseTursoUrl.String(), libsql.WithAuthToken(cfg.DatabaseTursoToken), libsql.WithSyncInterval(cfg.DatabaseTursoSyncTime))
	if err != nil {
		log.Fatalf("Failed to connect to turso database: failed to create connector: %v", err)
	}
	database.connector = connector

	database.db = sql.OpenDB(connector)
	database.Executor = sqlc.New(database.db)

	fmt.Printf("Created embedded replica at %s, syncing every %v\n", dbFile, cfg.DatabaseTursoSyncTime)
	return database
}
