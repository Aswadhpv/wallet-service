package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/yourusername/wallet-service/config"
	"github.com/yourusername/wallet-service/db"
	"github.com/yourusername/wallet-service/handler"
	"github.com/yourusername/wallet-service/repository"
	"github.com/yourusername/wallet-service/service"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	repo := repository.NewWalletRepo(database)
	svc := service.NewWalletService(repo)
	h := handler.NewWalletHandler(svc)

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallet", h.CreateOperation).Methods("POST")
	r.HandleFunc("/api/v1/wallets/{id}", h.GetBalance).Methods("GET")

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
