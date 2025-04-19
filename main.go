// @title        Wallet Service API
// @version      1.0
// @description  Simple wallet deposit/withdraw service.
// @host         localhost:8080
// @BasePath     /api/v1

package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    // Swagger middleware & generated docs
    httpSwagger "github.com/swaggo/http-swagger"
    _ "github.com/Aswadhpv/wallet-service/docs"

    "github.com/Aswadhpv/wallet-service/config"
    "github.com/Aswadhpv/wallet-service/db"
    "github.com/Aswadhpv/wallet-service/handler"
    "github.com/Aswadhpv/wallet-service/repository"
    "github.com/Aswadhpv/wallet-service/service"
)

func main() {
    // Load .env
    cfg := config.Load()

    // DB, repo, service, handler wiring
    database := db.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    repo := repository.NewWalletRepo(database)
    svc := service.NewWalletService(repo)
    h := handler.NewWalletHandler(svc)

    r := mux.NewRouter()

    // --- Your API routes under /api/v1 ---
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/wallet", h.CreateOperation).Methods("POST")
    api.HandleFunc("/wallets/{id}", h.GetBalance).Methods("GET")

    // --- Swagger UI ---
    // swagger.json is served automatically from /docs by swag
    // this will serve the interactive UI at http://localhost:8080/swagger/index.html
    r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    addr := fmt.Sprintf(":%s", cfg.AppPort)
    log.Printf("starting server on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
