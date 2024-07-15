package migration

import "database/sql"

type migration interface {
	Name() string
	Up(conn *sql.Tx) error
	Down(conn *sql.Tx) error
}

func getMigrations() []migration {
	return []migration{
		newMigration202407160214CreateTableUsers(),
		newMigration202407160215InsertTableUsers(),
	}
}


func checkDuplicateMigrationNames(migrations []migration) {
	nameSet := make(map[string]bool)
	for _, m := range migrations {
		if nameSet[m.Name()] {
			panic("Duplicate migration name found: " + m.Name())
		}
		nameSet[m.Name()] = true
	}
}

func Up(db *sql.DB) {
	// Get all migrations
	migrations := getMigrations()

	// Check for duplicate migration names
	checkDuplicateMigrationNames(migrations)

	// Create migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)

	if err != nil {
		panic(err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	for _, m := range migrations {
		// Check if migration has already been applied
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", m.Name()).Scan(&count)
		if err != nil {
			panic(err)
		}

		if count == 0 {
			// Apply migration
			if err := m.Up(tx); err != nil {
				panic(err)
			}

			// Record migration as applied
			_, err = tx.Exec("INSERT INTO migrations (name) VALUES ($1)", m.Name())
			if err != nil {
				panic(err)
			}

			println("Applied migration:", m.Name())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func Down(db *sql.DB) {
	// Get all migrations
	migrations := getMigrations()
	
	// Check for duplicate migration names
	checkDuplicateMigrationNames(migrations)

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	// Get the last applied migration
	var lastMigration string
	err = tx.QueryRow("SELECT name FROM migrations ORDER BY applied_at DESC LIMIT 1").Scan(&lastMigration)
	if err != nil {
		if err == sql.ErrNoRows {
			println("No migrations to revert")
			return
		}
		panic(err)
	}

	// Find the migration in our list
	var migrationToRevert migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Name() == lastMigration {
			migrationToRevert = migrations[i]
			break
		}
	}

	if migrationToRevert == nil {
		panic("Last applied migration not found in migration list")
	}

	// Revert the migration
	if err := migrationToRevert.Down(tx); err != nil {
		panic(err)
	}

	// Remove the migration record
	_, err = tx.Exec("DELETE FROM migrations WHERE name = $1", lastMigration)
	if err != nil {
		panic(err)
	}

	println("Reverted migration:", lastMigration)

	// Commit transaction
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func DownAll(db *sql.DB) {
	// Get all migrations
	migrations := getMigrations()
	
	// Check for duplicate migration names
	checkDuplicateMigrationNames(migrations)
	
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		// Check if migration has been applied
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", m.Name()).Scan(&count)
		if err != nil {
			panic(err)
		}

		if count > 0 {
			// Revert migration
			if err := m.Down(tx); err != nil {
				panic(err)
			}

			// Remove migration record
			_, err = tx.Exec("DELETE FROM migrations WHERE name = $1", m.Name())
			if err != nil {
				panic(err)
			}

			println("Reverted migration:", m.Name())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}
