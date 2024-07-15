package migration

import "database/sql"


type migration202407160214CreateTableUsers struct {
}

func newMigration202407160214CreateTableUsers() migration {
	return &migration202407160214CreateTableUsers{}
}


func (m *migration202407160214CreateTableUsers) Name() string {
	return "202407160214_create_table_users"
}

func (m *migration202407160214CreateTableUsers) Up(conn *sql.Tx) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func (m *migration202407160214CreateTableUsers) Down(conn *sql.Tx) error {
	_, err := conn.Exec(`
		DROP TABLE IF EXISTS users
	`)
	return err
}