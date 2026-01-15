package tools

import (
	"time"
)

// mockDB 是 DatabaseInterface 的模拟实现。
// 它为了测试和演示方便，使用内存映射 (map) 来存储数据，而不是真实的数据库。
type mockDB struct{}

// mockLoginDetails 模拟存储用户登录凭据的表。
// Key: 用户名, Value: LoginDetails 结构体。
var mockLoginDetails = map[string]LoginDetails{
	"alex": {
		Username:  "alex",
		AuthToken: "token_alex_123",
	},
	"bob": {
		Username:  "bob",
		AuthToken: "token_bob_456",
	},
}

// mockCoinDetails 模拟存储用户硬币余额的表。
// Key: 用户名, Value: CoinDetails 结构体。
var mockCoinDetails = map[string]CoinDetails{
	"alex": {
		Username: "alex",
		Balance:  1200,
	},
	"bob": {
		Username: "bob",
		Balance:  500,
	},
}

// GetUserLoginDetails 从内存 map 中查找用户的登录信息。
func (db *mockDB) GetUserLoginDetails(username string) *LoginDetails {
	// 模拟真实的数据库延迟，休眠 50 毫秒。
	time.Sleep(50 * time.Millisecond)

	var clientData = LoginDetails{}

	// 从 map 中查找 key 是否存在。
	clientData, ok := mockLoginDetails[username]
	if !ok {
		// 如果 key 不存在，返回 nil。
		return nil
	}

	// 返回找到的数据引用。
	return &clientData
}

// GetUserCoins 从内存 map 中查找用户的硬币余额信息。
func (db *mockDB) GetUserCoins(username string) *CoinDetails {
	// 模拟真实的数据库延迟，休眠 50 毫秒。
	time.Sleep(50 * time.Millisecond)

	var clientData = CoinDetails{}

	// 从 map 中查找 key 是否存在。
	clientData, ok := mockCoinDetails[username]
	if !ok {
		// 如果 key 不存在，返回 nil。
		return nil
	}

	// 返回找到的数据引用。
	return &clientData
}

// SetupDatabase 用于各种数据库初始化工作。
// 对于 mockDB，不需要做任何事情，直接返回 nil (无错误)。
func (db *mockDB) SetupDatabase() error {
	// No setup needed for mock database
	return nil
}
