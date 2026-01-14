package tools

import (
	"time"
)

type mockDB struct{}

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

func (db *mockDB) GetUserLoginDetails(username string) *LoginDetails {
	time.Sleep(50 * time.Millisecond) // Simulate database latency
	
	var clientData = LoginDetails{}
	clientData, ok := mockLoginDetails[username]
	if !ok {
		return nil
	}
	return &clientData
}

func (db *mockDB) GetUserCoins(username string) *CoinDetails {
	time.Sleep(50 * time.Millisecond) // Simulate database latency

	var clientData = CoinDetails{}
	clientData, ok := mockCoinDetails[username]
	if !ok {
		return nil
	}
	return &clientData
}

func (db *mockDB) SetupDatabase() error {
	// No setup needed for mock database
	return nil
}