package database

import (
	"context"
	"database/sql"
	"fmt"
	"time" // 新增导入

	// Import the MySQL driver anonymously to register it with database/sql

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhongjian-ready/goapi/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// MySQLDB is a wrapper around the standard sql.DB connection
type MySQLDB struct {
	db *sql.DB
}

// SetupDatabase configures the database connection using environment variables
func (d *MySQLDB) SetupDatabase(cfg *config.Config) error {
	var err error
	// Construct the Data Source Name (DSN) from config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
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

	// === 新增：配置数据库连接池 ===
	// 避免在高并发下创建过多连接导致数据库崩溃
	d.db.SetMaxOpenConns(25)                 // 最大打开连接数
	d.db.SetMaxIdleConns(25)                 // 最大空闲连接数
	d.db.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间

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
		password VARCHAR(255) NOT NULL
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
	// 为了演示，我们将密码设为 "123456" 的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 使用 IGNORE 关键字，如果数据已存在则忽略
	seedUserQuery := `INSERT IGNORE INTO users (id, username, password) VALUES (1, 'zhongjian', ?);`
	if _, err := d.db.Exec(seedUserQuery, string(hashedPassword)); err != nil {
		return err
	}

	seedCoinQuery := `INSERT IGNORE INTO coin_details (id, username, balance) VALUES (1, 'zhongjian', 1000);`
	if _, err := d.db.Exec(seedCoinQuery); err != nil {
		return err
	}

	return nil
}

// GetUserLoginDetails fetches the login credentials for a specific user
func (d *MySQLDB) GetUserLoginDetails(ctx context.Context, username string) (*LoginDetails, error) {
	var details LoginDetails
	query := "SELECT id, username, password FROM users WHERE username = ?"

	// Execute the query and scan the result into the struct
	// 使用 QueryRowContext 支持上下文取消
	err := d.db.QueryRowContext(ctx, query, username).Scan(&details.UserID, &details.Username, &details.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// 返回 nil 和 nil error 表示没找到
			return nil, nil
		}
		// 不要在这里打日志，把错误返回给调用者处理，避免双重日志
		return nil, err
	}

	return &details, nil
}

// GetUserCoins retrieves the coin balance for a user
func (d *MySQLDB) GetUserCoins(ctx context.Context, userid int) (*CoinDetails, error) {
	var details CoinDetails
	query := "SELECT id, username, balance FROM coin_details WHERE id = ?"

	err := d.db.QueryRowContext(ctx, query, userid).Scan(&details.UserID, &details.Username, &details.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &details, nil
}
