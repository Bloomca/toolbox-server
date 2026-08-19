package main

import (
	"database/sql"
	"fmt"
	"log"
)

type migration struct {
	name string
	up   string
}

var migrations = []migration{
	// Migrations are appended here in order.
}

func migrate(db *sql.DB) error {
	var currentVersion int
	if err := db.QueryRow("PRAGMA user_version").Scan(&currentVersion); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}

	if currentVersion < 0 || currentVersion > len(migrations) {
		return fmt.Errorf("database schema version %d is not supported", currentVersion)
	}

	for index := currentVersion; index < len(migrations); index++ {
		version := index + 1
		m := migrations[index]

		if err := applyMigration(db, version, m); err != nil {
			return fmt.Errorf("apply migration %d (%s): %w", version, m.name, err)
		}
		log.Printf("applied migration %d (%s)", version, m.name)
	}

	return nil
}

func applyMigration(db *sql.DB, version int, m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	if _, err := tx.Exec(m.up); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("execute SQL: %w", err)
	}

	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", version)); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record schema version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
