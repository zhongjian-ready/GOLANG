package tools

import (
	log "github.com/sirupsen/logrus"
)

// LoginDetails 存储用户的登录认证信息。
type LoginDetails struct {
	UserID    int    // 用户ID
	Username  string // 用户名
	AuthToken string // 认证令牌
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
type DatabaseInterface interface {
	GetUserLoginDetails(userid int) *LoginDetails // 根据用户ID获取登录信息
	GetUserCoins(userid int) *CoinDetails         // 根据用户ID获取硬币信息
	SetupDatabase() error                         // 初始化数据库连接或结构
	SetupSchema() error                           // (新增) 初始化数据库表结构（建表）
}

// NewDatabase 是一个工厂函数，用于创建和返回一个 DatabaseInterface 实例。
func NewDatabase() (*DatabaseInterface, error) {
	// 创建一个 MySQLDB 实例。
	var database DatabaseInterface = &MySQLDB{}

	// 初始化数据库连接
	err := database.SetupDatabase()
	if err != nil {
		log.Error("Failed to setup database:", err)
		return nil, err
	}

	// 返回数据库接口指针。
	return &database, nil
}
