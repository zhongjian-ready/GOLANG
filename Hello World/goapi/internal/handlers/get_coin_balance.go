package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/avukadin/goapi/api"
	"github.com/avukadin/goapi/internal/tools"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetCoinBalance(w http.ResponseWriter, r *http.Request) {
	var params = api.CoinBalanceParam{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error("Failed to decode request parameters:", err)
		api.RequestErrorHandler(w, err)
		return
	}
	
	var database *tools.DatabaseInterface
	database, err = tools.NewDatabase()

	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	var tokenDetails *tools.CoinDetails
	tokenDetails = (*database).GetUserCoins(params.UserName)
	
	if tokenDetails == nil {
		log.Error("Failed to get coin details for user:", params.UserName)
		api.RequestErrorHandler(w, err)
		return
	}

	var response = api.CoinBalanceResponse{
		Code : http.StatusOK,
		Balance : (*tokenDetails).Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	
	err = json.NewEncoder(w).Encode(response)
	
	if err != nil {
		log.Error("Failed to encode response:", err)
		api.InternalErrorHandler(w)
		return
	}


}