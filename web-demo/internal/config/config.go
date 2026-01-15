package config

import (
	"fmt"
	"os"
)

// Config 聚合了所有的应用配置
type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
}

// Load 从环境变量加载配置。
// 如果关键配置丢失，可以在这里返回错误，而不是在运行时崩溃。
func Load() (*Config, error) {
	cfg := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	// 简单的校验逻辑，确保生产环境不裸奔
	if cfg.JWTSecret == "" {
		// 开发环境可以给个默认值，但在生产环境这很危险
		// 这里为了演示方便，保持默认值逻辑，但打印警告
		cfg.JWTSecret = "default_secret_key"
	}

	// 可以在这里检查 DB 配置是否完整
	if cfg.DBUser == "" || cfg.DBHost == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}
