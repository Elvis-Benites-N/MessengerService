package db

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseLocal struct {
	db *sql.DB
}

// NewDatabase creates a new Database instance and establishes a connection to the database.
func NewDatabaseLocal() (*Database, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbName := os.Getenv("DB_NAME")

	// Check if all required environment variables are present
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbSSLMode == "" {
		return nil, errors.New("missing required environment variable(s)")
	}

	// Form the database connection string
	connectionDB := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " port=" + dbPort + " sslmode=" + dbSSLMode

	// Open the database connection
	db, err := sql.Open("postgres", connectionDB)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	// Check if the database already exists
	exists, err := checkDatabaseExists(db, dbName)
	if err != nil {
		return nil, err
	}

	if exists {
		// Database already exists, no need to create
		connectionDB = connectionDB + " dbname=" + dbName
		db, err = sql.Open("postgres", connectionDB)
		if err != nil {
			return nil, err
		}
	} else {

		// Database doesn't exist, create it
		createDBFilePath := filepath.Join("db", "migrate", "createDataBase.sql")
		err = executeSQLFile(db, createDBFilePath)
		if err != nil {
			return nil, err
		}

		connectionDB = connectionDB + " dbname=" + dbName
		db, err = sql.Open("postgres", connectionDB)
		if err != nil {
			return nil, err
		}
	}

	// Execute the SQL file to create tables
	filePath := filepath.Join("db", "migrate", "createTable.sql")
	err = executeSQLFile(db, filePath)
	if err != nil {
		return nil, err
	}

	// Uncomment the following lines if there is a need to execute a dropTable.sql file
	// dropTableFilePath := filepath.Join("db", "migrate", "dropTable.sql")
	// err = executeSQLFile(db, dropTableFilePath)
	// if err != nil {
	// 	return nil, err
	// }

	return &Database{db: db}, nil
}

// checkDatabaseExists checks if the database with the given name exists.
func checkDatabaseExists(db *sql.DB, dbName string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", dbName)
	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
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

// Close closes the database connection.
func (d *Database) CloseLocal() {
	d.db.Close()
}

// GetDB returns the underlying *sql.DB instance.
func (d *Database) GetDBLocal() *sql.DB {
	return d.db
}
