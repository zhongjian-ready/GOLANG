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

// SetupSchema creates necessary tables if they don't exist
func (d *MySQLDB) SetupSchema() error {
	// 1. Create users table
	// 对应 GetUserLoginDetails 中的查询
	createUsersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		auth_token VARCHAR(255) NOT NULL
	);`

	if _, err := d.db.Exec(createUsersQuery); err != nil {
		return err
	}

	// 2. Create coin_details table
	// 对应 GetUserCoins 中的查询
	createCoinsQuery := `
	CREATE TABLE IF NOT EXISTS coin_details (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		balance BIGINT NOT NULL
	);`

	if _, err := d.db.Exec(createCoinsQuery); err != nil {
		return err
	}

	// 3. (可选) 插入测试数据，为了方便演示
	// 使用 IGNORE 关键字，如果数据已存在则忽略，避免重复出错
	seedUserQuery := `INSERT IGNORE INTO users (id, username, auth_token) VALUES (1, 'zhongjian', '123456');`
	if _, err := d.db.Exec(seedUserQuery); err != nil {
		return err
	}

	seedCoinQuery := `INSERT IGNORE INTO coin_details (id, username, balance) VALUES (1, 'zhongjian', 1000);`
	if _, err := d.db.Exec(seedCoinQuery); err != nil {
		return err
	}

	return nil
}

// GetUserLoginDetails fetches the login credentials for a specific user
func (d *MySQLDB) GetUserLoginDetails(userid int) *LoginDetails {
	var details LoginDetails
	query := "SELECT id, username, auth_token FROM users WHERE id = ?"

	// Execute the query and scan the result into the struct
	err := d.db.QueryRow(query, userid).Scan(&details.UserID, &details.Username, &details.AuthToken)
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
func (d *MySQLDB) GetUserCoins(userid int) *CoinDetails {
	var details CoinDetails
	query := "SELECT id, username, balance FROM coin_details WHERE id = ?"

	err := d.db.QueryRow(query, userid).Scan(&details.UserID, &details.Username, &details.Balance)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &details
}
