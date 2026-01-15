package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/zhongjian-ready/goapi/internal/config"
	"github.com/zhongjian-ready/goapi/internal/database"
)

// MockDatabase 是 database.DatabaseInterface 接口的一个模拟实现（Mock）。
// 它基于 "github.com/stretchr/testify/mock" 包，允许我们在单元测试中
// 预设数据库的行为（例如指定某个查询返回什么数据，或者模拟数据库报错），
// 从而可以在不连接真实 MySQL 数据库的情况下测试业务逻辑。
type MockDatabase struct {
	mock.Mock
}

// GetUserLoginDetails 模拟获取用户登录信息的行为。
// 在测试中，可以使用 m.On("GetUserLoginDetails", ctx, username).Return(details, err) 来设置预期的返回值。
func (m *MockDatabase) GetUserLoginDetails(ctx context.Context, username string) (*database.LoginDetails, error) {
	// m.Called 记录了该方法的调用，并根据测试中的设定返回预设参数 (args)。
	args := m.Called(ctx, username)

	// args.Get(0) 获取第一个返回值（*database.LoginDetails）。
	// 如果设定为 nil，则手动转换 nil，避免类型断言 panic。
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// 否则进行类型断言并返回。
	return args.Get(0).(*database.LoginDetails), args.Error(1)
}

// GetUserCoins 模拟获取用户硬币余额的行为。
// 同样使用 m.On(...) 来设定特定 user_id 下的返回结果。
func (m *MockDatabase) GetUserCoins(ctx context.Context, userid int) (*database.CoinDetails, error) {
	args := m.Called(ctx, userid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.CoinDetails), args.Error(1)
}

// SetupDatabase 模拟数据库初始化的行为。
// 对于不需要实际测试数据库连接的场景，可以在测试中设定让它直接返回 nil (成功)。
func (m *MockDatabase) SetupDatabase(cfg *config.Config) error {
	args := m.Called(cfg)
	return args.Error(0)
}

// SetupSchema 模拟建表操作。
// 在单元测试中通常不需要真的建表，只需验证该方法是否被调用即可。
func (m *MockDatabase) SetupSchema() error {
	args := m.Called()
	return args.Error(0)
}
