package stores

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteDatabase(ctx context.Context, migrationUrl string, dbUrl string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		return nil, &StorageError{Op: "Unable to connect to sqlite db", Err: err}
	}
	runMigrations(db, migrationUrl)
	store := &Database{Db: db, stmts: make(map[string]*sql.Stmt)}
	err = store.PrepareStatements(ctx)
	if err != nil {
		return nil, &StorageError{Op: "Unable to connect to sqlite db", Err: err}
	}
	return store, nil
}

func runMigrations(db *sql.DB, p string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return &MigrationError{Op: "Unable to create sqlite driver", Err: err}
	}
	m, err := migrate.NewWithDatabaseInstance(
		p,
		"sqlite", driver,
	)
	if err != nil {
		return &MigrationError{Op: "Unable to create sqlite migrate instance", Err: err}
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return &MigrationError{Op: "Unable to run sqlite migrations", Err: err}
	}
	return nil
}

type StorageError struct {
	Op  string
	Err error
}

func (e *StorageError) Error() string { return e.Op + " " + e.Err.Error() }

func (e *StorageError) Unwrap() error { return e.Err }

type MigrationError struct {
	Op  string
	Err error
}

func (e *MigrationError) Error() string { return e.Op + " " + e.Err.Error() }

func (e *MigrationError) Unwrap() error { return e.Err }
