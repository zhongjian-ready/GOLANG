package tools

import (
	log "github.com/sirupsen/logrus"
)

type LoginDetails struct {
	Username  string
	AuthToken string
}

type CoinDetails struct {
	Username string
	Balance  int64
}

type DatabaseInterface interface {
	GetUserLoginDetails(username string) *LoginDetails
	GetUserCoins(username string) *CoinDetails
	SetupDatabase() error
}

func NewDatabase() (*DatabaseInterface, error) {
	var database DatabaseInterface = &mockDB{}
	err := database.SetupDatabase()
	if err != nil {
		log.Error("Failed to setup database:", err)
		return nil,  err
	}
	return &database, nil
}