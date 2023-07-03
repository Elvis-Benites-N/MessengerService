package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://root:password@db:5432/chatgo?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Add a delay to allow the database container to start up
	time.Sleep(5 * time.Second)

	// Execute the SQL file to create tables
	filePath := filepath.Join("db", "migrate", "createTable.sql")
	err = executeSQLFile(db, filePath)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

// executeSQLFile reads the SQL file and executes its content against the database.
func executeSQLFile(db *sql.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	_, err = db.Exec(string(readFile(file)))
	if err != nil {
		return err
	}

	return nil
}

// readFile reads the content of a file and returns it as a byte slice.
func readFile(file io.Reader) []byte {
	data, _ := io.ReadAll(file)
	return data
}
