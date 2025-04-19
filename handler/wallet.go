package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yourusername/wallet-service/service"
)

type WalletHandler struct {
	svc service.WalletService
}

func NewWalletHandler(s service.WalletService) *WalletHandler {
	return &WalletHandler{svc: s}
}

type operationRequest struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

func (h *WalletHandler) CreateOperation(w http.ResponseWriter, r *http.Request) {
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	var err error
	if req.OperationType == "DEPOSIT" {
		err = h.svc.Deposit(r.Context(), req.WalletID, req.Amount)
	} else {
		err = h.svc.Withdraw(r.Context(), req.WalletID, req.Amount)
	}
	if err != nil {
		if err == service.ErrInsufficientFunds {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	bal, err := h.svc.GetBalance(r.Context(), id)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]int64{"balance": bal})
}
