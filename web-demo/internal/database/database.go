package database

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zhongjian-ready/goapi/internal/config"
)

// LoginDetails 存储用户的登录认证信息。
type LoginDetails struct {
	UserID   int    // 用户ID
	Username string // 用户名
	Password string // 加密后的密码
}

// CoinDetails 存储用户的虚拟货币信息。
type CoinDetails struct {
	UserID   int    // 用户ID
	Username string // 用户名
	Balance  int64  // 余额
}

// DatabaseInterface 定义了数据访问层的统一接口。
// 任何实现了这些方法的结构体（如 MockDB, MySQLDB, PostgresDB）都可以作为 Database 使用。
// 这遵循了依赖倒置原则，方便测试和更换数据库实现。
// 2026-01-15 更新：所有查询方法增加 context 参数，支持取消和超时控制 (Standard Practice)。
type DatabaseInterface interface {
	GetUserLoginDetails(ctx context.Context, username string) (*LoginDetails, error) // 根据用户名获取登录信息
	GetUserCoins(ctx context.Context, userid int) (*CoinDetails, error)              // 根据用户ID获取硬币信息
	SetupDatabase(cfg *config.Config) error                                          // 初始化数据库连接或结构
	SetupSchema() error                                                              // (新增) 初始化数据库表结构（建表）
}

// NewDatabase 是一个工厂函数，用于创建和返回一个 DatabaseInterface 实例。
func NewDatabase(cfg *config.Config) (DatabaseInterface, error) {
	// 创建一个 MySQLDB 实例。
	var database DatabaseInterface = &MySQLDB{}

	// 初始化数据库连接
	err := database.SetupDatabase(cfg)
	if err != nil {
		log.Error("Failed to setup database:", err)
		return nil, err
	}

	// 返回数据库接口。
	return database, nil
}
