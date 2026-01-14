package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/avukadin/goapi/internal/handlers"
)

func main() {
	log.SetReportCaller(true)

	var r *chi.Mux = chi.NewRouter()

	handlers.Handler(r)

	fmt.Println("Start go api service...")

	fmt.Println(`
   ______  ____    ___    ____  ____
  / ____/ / __ \  /   |  / __ \/  _/
 / / __  / / / / / /| | / /_/ // /  
/ /_/ / / /_/ / / ___ |/ ____// /   
\____/  \____/ /_/  |_/_/   /___/   
`)

	err := http.ListenAndServe("localhost:8000", r)

	if err!= nil {
		log.Error(err)
	}
}
