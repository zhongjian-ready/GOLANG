package tools

import (
	"database/sql"
	"fmt"
	"os"

	// Import the MySQL driver anonymously to register it with database/sql

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// MySQLDB is a wrapper around the standard sql.DB connection
type MySQLDB struct {
	db *sql.DB
}

// SetupDatabase configures the database connection using environment variables
func (d *MySQLDB) SetupDatabase() error {
	var err error
	// Construct the Data Source Name (DSN) from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// Open a new connection pool to the database
	d.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection is actually alive
	if err = d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// GetUserLoginDetails fetches the login credentials for a specific user
func (d *MySQLDB) GetUserLoginDetails(username string) *LoginDetails {
	var details LoginDetails
	query := "SELECT username, auth_token FROM users WHERE username = ?"

	// Execute the query and scan the result into the struct
	err := d.db.QueryRow(query, username).Scan(&details.Username, &details.AuthToken)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err)
		}
		// Return nil if user not found or error occurred
		return nil
	}

	return &details
}

// GetUserCoins retrieves the coin balance for a user
func (d *MySQLDB) GetUserCoins(username string) *CoinDetails {
	var details CoinDetails
	query := "SELECT username, balance FROM coin_details WHERE username = ?"

	err := d.db.QueryRow(query, username).Scan(&details.Username, &details.Balance)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &details
}
