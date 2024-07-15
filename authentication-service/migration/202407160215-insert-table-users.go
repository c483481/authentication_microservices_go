package migration

import "database/sql"

type migration202407160215InsertTableUsers struct {
}

func newMigration202407160215InsertTableUsers() *migration202407160215InsertTableUsers {
	return &migration202407160215InsertTableUsers{}
}

func (m *migration202407160215InsertTableUsers) Name() string {
	return "202407160215_insert_table_users"
}

func (m *migration202407160215InsertTableUsers) Up(db *sql.Tx) error {
	_, err := db.Exec(`
		INSERT INTO users (email, first_name, last_name, password, active) VALUES
		('john.doe@example.com', 'John', 'Doe', '$2y$10$AsV1MFOLro7D5CM.9pAKJ.hAl32yXWxGGbtcG7/qLdWx/JaB8HVuK', true),
		('jane.doe@example.com', 'Jane', 'Doe', '$2y$10$MNVHEVBmUOyqTahbyTrVxOLNxxZU4I4p1gmnM5t/80EhbPsinpISa', true),
		('bob.smith@example.com', 'Bob', 'Smith', '$2y$10$MVmrBEke2ypC3XVKmRKUme5qdGFXXaCXH76tQ4T/lk.DNRc7sHw7.', true)
	`)
	return err
}

func (m *migration202407160215InsertTableUsers) Down(db *sql.Tx) error {
	_, err := db.Exec(`DELETE FROM users WHERE email IN ('john.doe@example.com', 'jane.doe@example.com', 'bob.smith@example.com')`)
	return err
}
